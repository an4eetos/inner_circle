// Package jobs wires background work on River (Postgres-backed queue).
// Feature modules register their workers and periodic jobs here.
package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

const (
	defaultQueueWorkers = 10
	heartbeatInterval   = time.Minute
)

// HeartbeatArgs is a periodic no-op proving the worker loop and queue are alive.
type HeartbeatArgs struct{}

// Kind identifies the job type in the river_job table.
func (HeartbeatArgs) Kind() string { return "heartbeat" }

// HeartbeatWorker handles HeartbeatArgs.
type HeartbeatWorker struct {
	river.WorkerDefaults[HeartbeatArgs]
	Logger *slog.Logger
}

// Work logs the beat.
func (w *HeartbeatWorker) Work(ctx context.Context, job *river.Job[HeartbeatArgs]) error {
	w.Logger.DebugContext(ctx, "heartbeat", "job_id", job.ID)
	return nil
}

// NewClient builds a River client with every worker and periodic job registered.
func NewClient(pool *pgxpool.Pool, log *slog.Logger) (*river.Client[pgx.Tx], error) {
	workers := river.NewWorkers()
	river.AddWorker(workers, &HeartbeatWorker{Logger: log})

	return river.NewClient(riverpgxv5.New(pool), &river.Config{
		Logger:  log,
		Workers: workers,
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: defaultQueueWorkers},
		},
		PeriodicJobs: []*river.PeriodicJob{
			river.NewPeriodicJob(
				river.PeriodicInterval(heartbeatInterval),
				func() (river.JobArgs, *river.InsertOpts) { return HeartbeatArgs{}, nil },
				&river.PeriodicJobOpts{RunOnStart: true},
			),
		},
	})
}
