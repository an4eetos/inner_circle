package httpx

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/an4eetos/inner_circle/backend/internal/platform/reqctx"
)

var ok = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

func TestRequestID(t *testing.T) {
	var seen string
	h := RequestID()(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		seen = reqctx.RequestID(r.Context())
	}))

	t.Run("reuses valid incoming id", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set(HeaderRequestID, "abc-123")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if seen != "abc-123" || w.Header().Get(HeaderRequestID) != "abc-123" {
			t.Fatalf("seen=%q header=%q", seen, w.Header().Get(HeaderRequestID))
		}
	})

	t.Run("replaces malformed id", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set(HeaderRequestID, "bad id with spaces\n")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if seen == "" || strings.Contains(seen, " ") || w.Header().Get(HeaderRequestID) != seen {
			t.Fatalf("expected generated id, got %q", seen)
		}
	})
}

func TestRecover(t *testing.T) {
	var logBuf bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logBuf, nil))
	h := Recover(log)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") }))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusInternalServerError || !strings.Contains(w.Body.String(), `"code":"internal"`) {
		t.Fatalf("got %d %s", w.Code, w.Body.String())
	}
	if !strings.Contains(logBuf.String(), "boom") {
		t.Fatalf("panic not logged: %s", logBuf.String())
	}
}

func TestAccessLog_UsesRoutePatternNotRawPath(t *testing.T) {
	var logBuf bytes.Buffer
	log := slog.New(slog.NewJSONHandler(&logBuf, nil))
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/things/{id}", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusTeapot) })
	h := AccessLog(log)(mux)

	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/v1/things/secret-id?q=private", nil))

	out := logBuf.String()
	if !strings.Contains(out, `"route":"GET /v1/things/{id}"`) || !strings.Contains(out, `"status":418`) {
		t.Fatalf("unexpected log: %s", out)
	}
	if strings.Contains(out, "secret-id") || strings.Contains(out, "private") {
		t.Fatalf("raw path or query leaked: %s", out)
	}
}

func TestCORS(t *testing.T) {
	h := CORS([]string{"http://app.test"})(ok)

	t.Run("allowed preflight", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodOptions, "/v1/healthz", nil)
		r.Header.Set("Origin", "http://app.test")
		r.Header.Set("Access-Control-Request-Method", "POST")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != http.StatusNoContent || w.Header().Get("Access-Control-Allow-Origin") != "http://app.test" {
			t.Fatalf("got %d %v", w.Code, w.Header())
		}
	})

	t.Run("unknown origin gets no CORS headers", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/v1/healthz", nil)
		r.Header.Set("Origin", "http://evil.test")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Fatalf("CORS header set for unknown origin")
		}
	})
}

func TestMaxBody(t *testing.T) {
	var readErr error
	h := MaxBody(4)(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		_, readErr = io.ReadAll(r.Body)
	}))

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/", strings.NewReader("too long")))
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("declared length over limit: got %d", w.Code)
	}

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("too long"))
	r.ContentLength = -1 // chunked: size unknown up front
	h.ServeHTTP(httptest.NewRecorder(), r)
	if readErr == nil {
		t.Fatal("expected read error past the limit")
	}
}

func TestRateLimit(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	h := RateLimit(RateLimitConfig{RPS: 1, Burst: 2, Now: func() time.Time { return now }})(ok)

	do := func(remote, xff string) int {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = remote
		if xff != "" {
			r.Header.Set("X-Forwarded-For", xff)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w.Code
	}

	if do("1.1.1.1:1", "") != 200 || do("1.1.1.1:2", "") != 200 {
		t.Fatal("burst should allow two requests")
	}
	if got := do("1.1.1.1:3", ""); got != http.StatusTooManyRequests {
		t.Fatalf("third request: got %d, want 429", got)
	}
	if do("2.2.2.2:1", "") != 200 {
		t.Fatal("other clients have their own bucket")
	}
	// Without TrustProxyHeaders, a spoofed X-Forwarded-For must not grant a fresh bucket.
	if got := do("1.1.1.1:4", "9.9.9.9"); got != http.StatusTooManyRequests {
		t.Fatalf("spoofed XFF bypassed limit: got %d", got)
	}
	now = now.Add(time.Second)
	if do("1.1.1.1:5", "") != 200 {
		t.Fatal("token should refill after 1s")
	}
}

func TestClientIP_TrustedProxyUsesLastHop(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RemoteAddr = "10.0.0.1:443"
	r.Header.Set("X-Forwarded-For", "6.6.6.6, 203.0.113.7")
	if got := clientIP(r, true); got != "203.0.113.7" {
		t.Fatalf("got %q", got)
	}
	if got := clientIP(r, false); got != "10.0.0.1" {
		t.Fatalf("got %q", got)
	}
}
