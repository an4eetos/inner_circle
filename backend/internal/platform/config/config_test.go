package config

import (
	"log/slog"
	"strings"
	"testing"
	"time"
)

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load(env(map[string]string{"DATABASE_URL": "postgres://x"}))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Env != EnvLocal || cfg.HTTPAddr != ":8080" || cfg.DBMaxConns != 10 ||
		cfg.LogLevel != slog.LevelInfo || cfg.LogFormat != "json" ||
		cfg.ShutdownTimeout != 15*time.Second || cfg.MaxBodyBytes != 1<<20 || cfg.OTLPEnabled {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestLoad_Overrides(t *testing.T) {
	cfg, err := Load(env(map[string]string{
		"APP_ENV":                     "staging",
		"DATABASE_URL":                "postgres://x",
		"LOG_LEVEL":                   "debug",
		"LOG_FORMAT":                  "text",
		"CORS_ALLOWED_ORIGINS":        " http://a.test , ,http://b.test",
		"RATE_LIMIT_RPS":              "2.5",
		"OTEL_EXPORTER_OTLP_ENDPOINT": "http://collector:4318",
	}))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Env != EnvStaging || cfg.LogLevel != slog.LevelDebug || cfg.LogFormat != "text" ||
		cfg.RateLimitRPS != 2.5 || !cfg.OTLPEnabled {
		t.Fatalf("overrides not applied: %+v", cfg)
	}
	if got := strings.Join(cfg.CORSAllowedOrigins, "|"); got != "http://a.test|http://b.test" {
		t.Fatalf("CORS origins = %q", got)
	}
}

func TestLoad_Errors(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want string
	}{
		{"missing db url", map[string]string{}, "DATABASE_URL"},
		{"bad env", map[string]string{"DATABASE_URL": "x", "APP_ENV": "qa"}, "APP_ENV"},
		{"auto migrate outside local", map[string]string{"DATABASE_URL": "x", "APP_ENV": "prod", "DB_AUTO_MIGRATE": "true"}, "DB_AUTO_MIGRATE"},
		{"bad int", map[string]string{"DATABASE_URL": "x", "DB_MAX_CONNS": "ten"}, "DB_MAX_CONNS"},
		{"bad duration", map[string]string{"DATABASE_URL": "x", "SHUTDOWN_TIMEOUT": "soon"}, "SHUTDOWN_TIMEOUT"},
		{"bad level", map[string]string{"DATABASE_URL": "x", "LOG_LEVEL": "loud"}, "LOG_LEVEL"},
		{"bad format", map[string]string{"DATABASE_URL": "x", "LOG_FORMAT": "xml"}, "LOG_FORMAT"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Load(env(tc.env))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want error mentioning %s, got %v", tc.want, err)
			}
		})
	}
}
