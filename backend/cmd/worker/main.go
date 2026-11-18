// Command worker runs River background jobs and periodic tasks.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/an4eetos/inner_circle/backend/internal/jobs"
	"github.com/an4eetos/inner_circle/backend/internal/platform/buildinfo"
	"github.com/an4eetos/inner_circle/backend/internal/platform/config"
	"github.com/an4eetos/inner_circle/backend/internal/platform/db"
	"github.com/an4eetos/inner_circle/backend/internal/platform/logging"
	"github.com/an4eetos/inner_circle/backend/internal/platform/telemetry"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "worker:", err)
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

	shutdownTelemetry, err := telemetry.Setup(ctx, cfg.OTLPEnabled, "inner-circle-worker", version)
	if err != nil {
		return err
	}

	pool, err := db.Connect(ctx, cfg.DatabaseURL, cfg.DBMaxConns)
	if err != nil {
		return err
	}
	defer pool.Close()

	client, err := jobs.NewClient(pool, log)
	if err != nil {
		return fmt.Errorf("river client: %w", err)
	}
	// Start with a context that outlives the signal so in-flight jobs can finish during Stop.
	if err := client.Start(context.WithoutCancel(ctx)); err != nil {
		return fmt.Errorf("river start: %w", err)
	}
	log.Info("worker started", "env", cfg.Env, "version", version)

	<-ctx.Done()
	log.Info("worker shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := client.Stop(shutdownCtx); err != nil {
		log.Warn("graceful stop timed out; cancelling running jobs", "error", err)
		if err := client.StopAndCancel(context.Background()); err != nil {
			log.Error("river stop and cancel", "error", err)
		}
	}
	if err := shutdownTelemetry(shutdownCtx); err != nil {
		log.Error("telemetry shutdown", "error", err)
	}
	return nil
}
