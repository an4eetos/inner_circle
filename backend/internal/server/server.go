// Package server assembles the HTTP handler: generated routes, module handlers and middleware.
package server

import (
	"errors"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/an4eetos/inner_circle/backend/internal/api"
	"github.com/an4eetos/inner_circle/backend/internal/health"
	"github.com/an4eetos/inner_circle/backend/internal/platform/config"
	"github.com/an4eetos/inner_circle/backend/internal/platform/httpx"
	"github.com/an4eetos/inner_circle/backend/internal/platform/problem"
)

// Deps are the collaborators the HTTP layer needs.
type Deps struct {
	Config config.Config
	Logger *slog.Logger
	Health *health.Handler
}

// handlers composes the per-module handlers into api.StrictServerInterface.
// Each new module adds its handler as an embedded field here.
type handlers struct {
	*health.Handler
}

var _ api.StrictServerInterface = handlers{}

// New returns the root HTTP handler.
func New(d Deps) http.Handler {
	mux := http.NewServeMux()

	strict := api.NewStrictHandlerWithOptions(handlers{Handler: d.Health}, nil, api.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  requestError,
		ResponseErrorHandlerFunc: responseError(d.Logger),
	})
	api.HandlerWithOptions(strict, api.StdHTTPServerOptions{
		BaseRouter:       mux,
		ErrorHandlerFunc: requestError,
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		problem.Write(w, r, http.StatusNotFound, problem.CodeNotFound, "")
	})

	return httpx.Chain(mux,
		httpx.Recover(d.Logger),
		httpx.RequestID(),
		func(h http.Handler) http.Handler { return otelhttp.NewHandler(h, "http.server") },
		httpx.CORS(d.Config.CORSAllowedOrigins),
		httpx.RateLimit(httpx.RateLimitConfig{
			RPS:               d.Config.RateLimitRPS,
			Burst:             d.Config.RateLimitBurst,
			TrustProxyHeaders: d.Config.TrustProxyHeaders,
		}),
		httpx.MaxBody(d.Config.MaxBodyBytes),
		httpx.AccessLog(d.Logger, "/v1/healthz", "/v1/readyz"),
	)
}

// requestError answers malformed requests (bad params, unreadable bodies).
func requestError(w http.ResponseWriter, r *http.Request, err error) {
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		problem.Write(w, r, http.StatusRequestEntityTooLarge, problem.CodePayloadTooLarge, "")
		return
	}
	problem.Write(w, r, http.StatusBadRequest, problem.CodeBadRequest, err.Error())
}

// responseError answers handler errors. Details are logged, never returned to the client.
func responseError(log *slog.Logger) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		log.ErrorContext(r.Context(), "handler error", "error", err)
		problem.Write(w, r, http.StatusInternalServerError, problem.CodeInternal, "")
	}
}
