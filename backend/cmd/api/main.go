// Command api serves the Inner Circle HTTP API.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/an4eetos/inner_circle/backend/internal/health"
	"github.com/an4eetos/inner_circle/backend/internal/platform/buildinfo"
	"github.com/an4eetos/inner_circle/backend/internal/platform/config"
	"github.com/an4eetos/inner_circle/backend/internal/platform/db"
	"github.com/an4eetos/inner_circle/backend/internal/platform/db/platformdb"
	"github.com/an4eetos/inner_circle/backend/internal/platform/logging"
	"github.com/an4eetos/inner_circle/backend/internal/platform/telemetry"
	"github.com/an4eetos/inner_circle/backend/internal/server"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 60 * time.Second
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "api:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		return err
	}
	log := logging.New(os.Stdout, cfg.LogLevel, cfg.LogFormat)
	slog.SetDefault(log)
	version := buildinfo.Get()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownTelemetry, err := telemetry.Setup(ctx, cfg.OTLPEnabled, "inner-circle-api", version)
	if err != nil {
		return err
	}

	pool, err := db.Connect(ctx, cfg.DatabaseURL, cfg.DBMaxConns)
	if err != nil {
		return err
	}
	defer pool.Close()

	if cfg.DBAutoMigrate {
		if err := db.Migrate(ctx, pool, log); err != nil {
			return err
		}
	}

	srv := &http.Server{
		Addr: cfg.HTTPAddr,
		Handler: server.New(server.Deps{
			Config: cfg,
			Logger: log,
			Health: health.New(platformdb.New(pool), version, log),
		}),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelWarn),
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("api listening", "addr", cfg.HTTPAddr, "env", cfg.Env, "version", version)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("listen: %w", err)
	case <-ctx.Done():
	}

	log.Info("api shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("http shutdown", "error", err)
	}
	if err := shutdownTelemetry(shutdownCtx); err != nil {
		log.Error("telemetry shutdown", "error", err)
	}
	return nil
}
