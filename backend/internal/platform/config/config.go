// Package config loads process configuration from environment variables.
// Nothing environment-specific (hosts, ports, credentials) is hard-coded elsewhere.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

// Env is the deployment environment.
type Env string

// Supported environments.
const (
	EnvLocal   Env = "local"
	EnvStaging Env = "staging"
	EnvProd    Env = "prod"
)

const maxDBConns = 1000

// Config is the runtime configuration shared by the api, worker and migrate binaries.
type Config struct {
	Env                Env
	HTTPAddr           string
	DatabaseURL        string
	DBMaxConns         int32
	DBAutoMigrate      bool
	LogLevel           slog.Level
	LogFormat          string
	CORSAllowedOrigins []string
	RateLimitRPS       float64
	RateLimitBurst     int
	TrustProxyHeaders  bool
	MaxBodyBytes       int64
	ShutdownTimeout    time.Duration
	OTLPEnabled        bool
}

// Load reads configuration through getenv (pass os.Getenv in production, a map lookup in tests).
func Load(getenv func(string) string) (Config, error) {
	p := parser{getenv: getenv}
	cfg := Config{
		Env:                Env(p.str("APP_ENV", string(EnvLocal))),
		HTTPAddr:           p.str("HTTP_ADDR", ":8080"),
		DatabaseURL:        p.str("DATABASE_URL", ""),
		DBAutoMigrate:      p.boolean("DB_AUTO_MIGRATE", false),
		LogLevel:           p.level("LOG_LEVEL", slog.LevelInfo),
		LogFormat:          p.str("LOG_FORMAT", "json"),
		CORSAllowedOrigins: p.list("CORS_ALLOWED_ORIGINS"),
		RateLimitRPS:       p.float("RATE_LIMIT_RPS", 20),
		RateLimitBurst:     p.integer("RATE_LIMIT_BURST", 40),
		TrustProxyHeaders:  p.boolean("TRUST_PROXY_HEADERS", false),
		MaxBodyBytes:       int64(p.integer("MAX_BODY_BYTES", 1<<20)),
		ShutdownTimeout:    p.duration("SHUTDOWN_TIMEOUT", 15*time.Second),
		OTLPEnabled: getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" ||
			getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT") != "",
	}

	maxConns := p.integer("DB_MAX_CONNS", 10)
	errs := p.errs
	switch cfg.Env {
	case EnvLocal, EnvStaging, EnvProd:
	default:
		errs = append(errs, fmt.Errorf("APP_ENV: must be local, staging or prod, got %q", cfg.Env))
	}
	if cfg.DatabaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL: required"))
	}
	if cfg.DBAutoMigrate && cfg.Env != EnvLocal {
		errs = append(errs, errors.New("DB_AUTO_MIGRATE: only allowed when APP_ENV=local; run the migrate binary in deploys"))
	}
	if cfg.LogFormat != "json" && cfg.LogFormat != "text" {
		errs = append(errs, fmt.Errorf("LOG_FORMAT: must be json or text, got %q", cfg.LogFormat))
	}
	if maxConns < 1 || maxConns > maxDBConns {
		errs = append(errs, fmt.Errorf("DB_MAX_CONNS: must be between 1 and %d", maxDBConns))
	} else {
		cfg.DBMaxConns = int32(maxConns)
	}
	if cfg.RateLimitRPS <= 0 || cfg.RateLimitBurst < 1 {
		errs = append(errs, errors.New("RATE_LIMIT_RPS and RATE_LIMIT_BURST: must be positive"))
	}
	if cfg.MaxBodyBytes < 1 {
		errs = append(errs, errors.New("MAX_BODY_BYTES: must be >= 1"))
	}
	if err := errors.Join(errs...); err != nil {
		return Config{}, fmt.Errorf("config: %w", err)
	}
	return cfg, nil
}

type parser struct {
	getenv func(string) string
	errs   []error
}

func (p *parser) str(key, def string) string {
	if v := strings.TrimSpace(p.getenv(key)); v != "" {
		return v
	}
	return def
}

func (p *parser) integer(key string, def int) int {
	v := strings.TrimSpace(p.getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		p.errs = append(p.errs, fmt.Errorf("%s: not an integer: %q", key, v))
		return def
	}
	return n
}

func (p *parser) float(key string, def float64) float64 {
	v := strings.TrimSpace(p.getenv(key))
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		p.errs = append(p.errs, fmt.Errorf("%s: not a number: %q", key, v))
		return def
	}
	return f
}

func (p *parser) boolean(key string, def bool) bool {
	v := strings.TrimSpace(p.getenv(key))
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		p.errs = append(p.errs, fmt.Errorf("%s: not a boolean: %q", key, v))
		return def
	}
	return b
}

func (p *parser) duration(key string, def time.Duration) time.Duration {
	v := strings.TrimSpace(p.getenv(key))
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		p.errs = append(p.errs, fmt.Errorf("%s: not a duration: %q", key, v))
		return def
	}
	return d
}

func (p *parser) level(key string, def slog.Level) slog.Level {
	v := strings.TrimSpace(p.getenv(key))
	if v == "" {
		return def
	}
	var l slog.Level
	if err := l.UnmarshalText([]byte(v)); err != nil {
		p.errs = append(p.errs, fmt.Errorf("%s: unknown level %q", key, v))
		return def
	}
	return l
}

func (p *parser) list(key string) []string {
	var out []string
	for _, item := range strings.Split(p.getenv(key), ",") {
		if item = strings.TrimSpace(item); item != "" {
			out = append(out, item)
		}
	}
	return out
}
