// Package db opens the PostgreSQL pool and applies migrations.
package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"

	"github.com/an4eetos/inner_circle/backend/migrations"
)

const connectTimeout = 10 * time.Second

// Connect opens a pool and verifies it with a ping.
func Connect(ctx context.Context, url string, maxConns int32) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("db: parse DATABASE_URL: %w", err)
	}
	cfg.MaxConns = maxConns

	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("db: open pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: ping: %w", err)
	}
	return pool, nil
}

// Migrate applies all pending goose migrations, then River's queue schema. Idempotent.
func Migrate(ctx context.Context, pool *pgxpool.Pool, log *slog.Logger) error {
	provider, closeDB, err := newGooseProvider(pool)
	if err != nil {
		return err
	}
	defer closeDB()

	results, err := provider.Up(ctx)
	if err != nil {
		return fmt.Errorf("db: goose up: %w", err)
	}
	for _, r := range results {
		log.InfoContext(ctx, "migration applied", "version", r.Source.Version, "duration_ms", r.Duration.Milliseconds())
	}

	migrator, err := rivermigrate.New(riverpgxv5.New(pool), nil)
	if err != nil {
		return fmt.Errorf("db: river migrator: %w", err)
	}
	res, err := migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	if err != nil {
		return fmt.Errorf("db: river migrate: %w", err)
	}
	for _, v := range res.Versions {
		log.InfoContext(ctx, "river migration applied", "version", v.Version)
	}
	return nil
}

// MigrateDownOne rolls back the most recent goose migration (River's schema is left alone).
func MigrateDownOne(ctx context.Context, pool *pgxpool.Pool) (int64, error) {
	provider, closeDB, err := newGooseProvider(pool)
	if err != nil {
		return 0, err
	}
	defer closeDB()
	r, err := provider.Down(ctx)
	if err != nil {
		return 0, fmt.Errorf("db: goose down: %w", err)
	}
	return r.Source.Version, nil
}

// MigrationStatus lists goose migrations and whether each is applied.
func MigrationStatus(ctx context.Context, pool *pgxpool.Pool) ([]*goose.MigrationStatus, error) {
	provider, closeDB, err := newGooseProvider(pool)
	if err != nil {
		return nil, err
	}
	defer closeDB()
	return provider.Status(ctx)
}

func newGooseProvider(pool *pgxpool.Pool) (*goose.Provider, func(), error) {
	sqlDB := stdlib.OpenDBFromPool(pool)
	provider, err := goose.NewProvider(goose.DialectPostgres, sqlDB, migrations.FS)
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("db: goose provider: %w", err)
	}
	return provider, func() { _ = sqlDB.Close() }, nil
}
