// Command migrate applies database migrations: goose SQL files plus River's queue schema.
//
//	migrate up       apply everything pending (default)
//	migrate down     roll back the latest goose migration
//	migrate status   list goose migrations
//
// Deploys run `migrate up` as a separate step before rolling out api and worker.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/an4eetos/inner_circle/backend/internal/platform/config"
	"github.com/an4eetos/inner_circle/backend/internal/platform/db"
	"github.com/an4eetos/inner_circle/backend/internal/platform/logging"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "migrate:", err)
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.Connect(ctx, cfg.DatabaseURL, 2)
	if err != nil {
		return err
	}
	defer pool.Close()

	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "up":
		if err := db.Migrate(ctx, pool, log); err != nil {
			return err
		}
		log.Info("migrations up to date")
	case "down":
		v, err := db.MigrateDownOne(ctx, pool)
		if err != nil {
			return err
		}
		log.Info("rolled back", "version", v)
	case "status":
		statuses, err := db.MigrationStatus(ctx, pool)
		if err != nil {
			return err
		}
		for _, s := range statuses {
			fmt.Printf("%-8s %05d  %s\n", s.State, s.Source.Version, s.Source.Path)
		}
	default:
		return fmt.Errorf("unknown command %q (want up, down or status)", cmd)
	}
	return nil
}
