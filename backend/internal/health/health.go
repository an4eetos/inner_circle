// Package health implements the liveness and readiness endpoints.
package health

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/an4eetos/inner_circle/backend/internal/api"
	"github.com/an4eetos/inner_circle/backend/internal/platform/problem"
)

const readinessTimeout = 2 * time.Second

// DB is the slice of the platform queries the readiness check needs.
type DB interface {
	DatabaseNow(ctx context.Context) (time.Time, error)
}

// Handler serves GET /v1/healthz and GET /v1/readyz.
type Handler struct {
	db      DB
	version string
	log     *slog.Logger
	now     func() time.Time
}

// New returns a health handler.
func New(db DB, version string, log *slog.Logger) *Handler {
	return &Handler{db: db, version: version, log: log, now: time.Now}
}

// GetHealthz reports liveness without touching dependencies.
func (h *Handler) GetHealthz(context.Context, api.GetHealthzRequestObject) (api.GetHealthzResponseObject, error) {
	return api.GetHealthz200JSONResponse{
		Status:  api.HealthStatusOk,
		Version: h.version,
		Time:    h.now().UTC(),
	}, nil
}

// GetReadyz reports whether the database answers a query.
func (h *Handler) GetReadyz(ctx context.Context, _ api.GetReadyzRequestObject) (api.GetReadyzResponseObject, error) {
	qctx, cancel := context.WithTimeout(ctx, readinessTimeout)
	defer cancel()
	if _, err := h.db.DatabaseNow(qctx); err != nil {
		h.log.WarnContext(ctx, "readiness: database check failed", "error", err)
		return api.GetReadyz503ApplicationProblemPlusJSONResponse(
			problem.New(ctx, http.StatusServiceUnavailable, problem.CodeUnavailable, "database unavailable"),
		), nil
	}
	return api.GetReadyz200JSONResponse{
		Status: api.ReadinessStatusOk,
		Checks: api.ReadinessChecks{Database: api.ReadinessChecksDatabaseOk},
	}, nil
}
