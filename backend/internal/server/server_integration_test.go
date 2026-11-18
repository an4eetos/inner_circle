package server_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/an4eetos/inner_circle/backend/internal/api"
	"github.com/an4eetos/inner_circle/backend/internal/health"
	"github.com/an4eetos/inner_circle/backend/internal/platform/config"
	"github.com/an4eetos/inner_circle/backend/internal/platform/db"
	"github.com/an4eetos/inner_circle/backend/internal/platform/db/platformdb"
	"github.com/an4eetos/inner_circle/backend/internal/platform/problem"
	"github.com/an4eetos/inner_circle/backend/internal/server"
)

// postgresImage must match infra/docker-compose.yml.
const postgresImage = "postgres:17-alpine"

func TestServer_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test: needs Docker")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	ctr, err := postgres.Run(ctx, postgresImage,
		postgres.WithDatabase("inner_circle"),
		postgres.WithUsername("inner_circle"),
		postgres.WithPassword("inner_circle"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("start postgres: %v", err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(context.Background()) })

	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("dsn: %v", err)
	}
	pool, err := db.Connect(ctx, dsn, 4)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	log := slog.New(slog.DiscardHandler)

	// Migrations apply cleanly and are idempotent.
	for i := range 2 {
		if err := db.Migrate(ctx, pool, log); err != nil {
			t.Fatalf("migrate run %d: %v", i+1, err)
		}
	}
	var hasCitext, hasRiver bool
	if err := pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'citext')`).Scan(&hasCitext); err != nil || !hasCitext {
		t.Fatalf("citext extension missing (err=%v)", err)
	}
	if err := pool.QueryRow(ctx, `SELECT to_regclass('public.river_job') IS NOT NULL`).Scan(&hasRiver); err != nil || !hasRiver {
		t.Fatalf("river tables missing (err=%v)", err)
	}

	cfg := config.Config{RateLimitRPS: 1000, RateLimitBurst: 1000, MaxBodyBytes: 1 << 20}
	srv := httptest.NewServer(server.New(server.Deps{
		Config: cfg,
		Logger: log,
		Health: health.New(platformdb.New(pool), "test-version", log),
	}))
	t.Cleanup(srv.Close)

	t.Run("healthz", func(t *testing.T) {
		var h api.Health
		hdr := get(t, srv.URL+"/v1/healthz", http.StatusOK, "application/json", &h)
		if h.Status != api.HealthStatusOk || h.Version != "test-version" || h.Time.IsZero() {
			t.Fatalf("unexpected body: %+v", h)
		}
		if hdr.Get("X-Request-ID") == "" {
			t.Fatal("missing X-Request-ID")
		}
	})

	t.Run("readyz", func(t *testing.T) {
		var r api.Readiness
		get(t, srv.URL+"/v1/readyz", http.StatusOK, "application/json", &r)
		if r.Status != api.ReadinessStatusOk || r.Checks.Database != api.ReadinessChecksDatabaseOk {
			t.Fatalf("unexpected body: %+v", r)
		}
	})

	t.Run("unknown route is a problem+json 404", func(t *testing.T) {
		var p api.Problem
		hdr := get(t, srv.URL+"/v1/nope", http.StatusNotFound, problem.ContentType, &p)
		if p.Code != problem.CodeNotFound || p.RequestId == nil || *p.RequestId != hdr.Get("X-Request-ID") {
			t.Fatalf("unexpected problem: %+v", p)
		}
	})

	t.Run("readyz is 503 when the database is gone", func(t *testing.T) {
		pool.Close()
		var p api.Problem
		get(t, srv.URL+"/v1/readyz", http.StatusServiceUnavailable, problem.ContentType, &p)
		if p.Code != problem.CodeUnavailable {
			t.Fatalf("unexpected problem: %+v", p)
		}
	})
}

// get fetches url, checks status and content type, decodes JSON into `into`, and returns the headers.
func get(t *testing.T, url string, wantStatus int, wantContentType string, into any) http.Header {
	t.Helper()
	res, err := http.Get(url) //nolint:gosec,noctx // test-only request to a local httptest server
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != wantStatus {
		t.Fatalf("GET %s: status %d, want %d; body %s", url, res.StatusCode, wantStatus, body)
	}
	if ct := res.Header.Get("Content-Type"); ct != wantContentType {
		t.Fatalf("GET %s: content-type %q, want %q", url, ct, wantContentType)
	}
	if err := json.Unmarshal(body, into); err != nil {
		t.Fatalf("decode %s: %v", body, err)
	}
	return res.Header
}
