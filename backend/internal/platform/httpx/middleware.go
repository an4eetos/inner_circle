// Package httpx holds HTTP middleware shared by every API route.
package httpx

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"regexp"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"

	"github.com/an4eetos/inner_circle/backend/internal/platform/problem"
	"github.com/an4eetos/inner_circle/backend/internal/platform/reqctx"
)

// HeaderRequestID is accepted from clients (if well-formed) and always echoed on responses.
const HeaderRequestID = "X-Request-ID"

// Middleware wraps a handler.
type Middleware func(http.Handler) http.Handler

// Chain applies middleware so the first one listed is the outermost.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

var validRequestID = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

// RequestID reuses a well-formed incoming X-Request-ID or generates a UUIDv7.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(HeaderRequestID)
			if !validRequestID.MatchString(id) {
				id = newID()
			}
			w.Header().Set(HeaderRequestID, id)
			next.ServeHTTP(w, r.WithContext(reqctx.WithRequestID(r.Context(), id)))
		})
	}
}

func newID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.NewString()
	}
	return id.String()
}

// Recover turns a panic into a logged 500 problem response.
func Recover(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				rec := recover()
				if rec == nil {
					return
				}
				if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
					panic(rec)
				}
				log.ErrorContext(r.Context(), "panic serving request",
					"panic", fmt.Sprint(rec), "stack", string(debug.Stack()))
				problem.Write(w, r, http.StatusInternalServerError, problem.CodeInternal, "")
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// AccessLog logs one line per request. It logs the matched route pattern rather than the raw
// path (paths will carry IDs) and never logs query strings or bodies. It must sit directly
// outside the ServeMux so it sees r.Pattern after routing.
func AccessLog(log *slog.Logger, quietPaths ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)

			route := r.Pattern
			if route == "" {
				route = "unmatched"
			}
			level := slog.LevelInfo
			switch {
			case rec.status >= http.StatusInternalServerError:
				level = slog.LevelError
			case slices.Contains(quietPaths, r.URL.Path):
				level = slog.LevelDebug
			default:
			}
			log.LogAttrs(r.Context(), level, "http request",
				slog.String("method", r.Method),
				slog.String("route", route),
				slog.Int("status", rec.status),
				slog.Int64("bytes", rec.bytes),
				slog.Float64("duration_ms", float64(time.Since(start).Microseconds())/1000),
			)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status      int
	bytes       int64
	wroteHeader bool
}

func (s *statusRecorder) WriteHeader(code int) {
	if !s.wroteHeader {
		s.status = code
		s.wroteHeader = true
	}
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	s.wroteHeader = true
	n, err := s.ResponseWriter.Write(b)
	s.bytes += int64(n)
	return n, err
}

// Unwrap lets http.ResponseController reach the underlying writer.
func (s *statusRecorder) Unwrap() http.ResponseWriter { return s.ResponseWriter }

// CORS allows browser calls from the configured origins only. An empty list disables CORS.
func CORS(allowedOrigins []string) Middleware {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if _, ok := allowed[origin]; !ok || origin == "" {
				next.ServeHTTP(w, r)
				return
			}
			h := w.Header()
			h.Add("Vary", "Origin")
			h.Set("Access-Control-Allow-Origin", origin)
			h.Set("Access-Control-Expose-Headers", HeaderRequestID)
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
				h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Idempotency-Key, X-Request-ID")
				h.Set("Access-Control-Max-Age", "600")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// MaxBody caps request body size; handlers reading past it get an error and the
// strict handler's request-error path answers 413/400.
func MaxBody(limit int64) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > limit {
				problem.Write(w, r, http.StatusRequestEntityTooLarge, problem.CodePayloadTooLarge, "")
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitConfig configures the per-client token bucket.
type RateLimitConfig struct {
	RPS               float64
	Burst             int
	TrustProxyHeaders bool // use the last X-Forwarded-For hop; only behind a proxy that sets it
	IdleTTL           time.Duration
	Now               func() time.Time
}

// RateLimit is a per-client-IP token bucket. It is a coarse abuse guard for a single
// instance; per-route and per-user limits (auth, invites, nudges) are added with those
// features. Rejected requests get 429 with Retry-After.
func RateLimit(cfg RateLimitConfig) Middleware {
	if cfg.IdleTTL == 0 {
		cfg.IdleTTL = 10 * time.Minute
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	l := &ipLimiter{
		limit:   rate.Limit(cfg.RPS),
		burst:   cfg.Burst,
		ttl:     cfg.IdleTTL,
		clients: make(map[string]*client),
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !l.allow(clientIP(r, cfg.TrustProxyHeaders), cfg.Now()) {
				w.Header().Set("Retry-After", strconv.Itoa(1))
				problem.Write(w, r, http.StatusTooManyRequests, problem.CodeRateLimited, "")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type client struct {
	lim  *rate.Limiter
	seen time.Time
}

type ipLimiter struct {
	mu        sync.Mutex
	limit     rate.Limit
	burst     int
	ttl       time.Duration
	clients   map[string]*client
	lastSweep time.Time
}

func (l *ipLimiter) allow(ip string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if now.Sub(l.lastSweep) > l.ttl {
		for k, c := range l.clients {
			if now.Sub(c.seen) > l.ttl {
				delete(l.clients, k)
			}
		}
		l.lastSweep = now
	}
	c, ok := l.clients[ip]
	if !ok {
		c = &client{lim: rate.NewLimiter(l.limit, l.burst)}
		l.clients[ip] = c
	}
	c.seen = now
	return c.lim.AllowN(now, 1)
}

func clientIP(r *http.Request, trustProxy bool) string {
	if trustProxy {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			parts := strings.Split(xff, ",")
			if ip := strings.TrimSpace(parts[len(parts)-1]); ip != "" {
				return ip
			}
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
