// Package logging builds the process logger: structured slog output with request/trace IDs
// attached and health or credential fields scrubbed before they reach any sink.
package logging

import (
	"context"
	"io"
	"log/slog"
	"strings"

	"go.opentelemetry.io/otel/trace"

	"github.com/an4eetos/inner_circle/backend/internal/platform/reqctx"
)

// Redacted replaces the value of any sensitive attribute.
const Redacted = "[REDACTED]"

// sensitiveKeys are attribute keys whose values never reach logs (case-insensitive).
// Health content (meds, meals, results, notes) and credentials are covered here; any
// key starting with "health_" is also scrubbed. See docs/architecture/00-system-overview.md.
var sensitiveKeys = map[string]struct{}{
	"authorization": {}, "cookie": {}, "set-cookie": {},
	"password": {}, "otp": {}, "secret": {},
	"token": {}, "access_token": {}, "refresh_token": {}, "id_token": {},
	"email": {}, "phone": {},
	"note": {}, "notes": {}, "comment": {}, "body": {},
	"meal": {}, "meals": {}, "medication": {}, "meds": {}, "dose": {},
	"result": {}, "results": {}, "lab": {}, "lab_value": {}, "lab_values": {}, "document": {},
}

// IsSensitive reports whether an attribute key must be redacted.
func IsSensitive(key string) bool {
	k := strings.ToLower(key)
	if _, ok := sensitiveKeys[k]; ok {
		return true
	}
	return strings.HasPrefix(k, "health_")
}

// New returns a logger writing to w in "json" or "text" format.
func New(w io.Writer, level slog.Level, format string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: level}
	var h slog.Handler
	if format == "text" {
		h = slog.NewTextHandler(w, opts)
	} else {
		h = slog.NewJSONHandler(w, opts)
	}
	return slog.New(&handler{next: h})
}

// handler scrubs sensitive attributes and adds request_id / trace_id from the context.
type handler struct {
	next slog.Handler
}

func (h *handler) Enabled(ctx context.Context, l slog.Level) bool { return h.next.Enabled(ctx, l) }

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	out := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	r.Attrs(func(a slog.Attr) bool {
		out.AddAttrs(scrub(a))
		return true
	})
	if id := reqctx.RequestID(ctx); id != "" {
		out.AddAttrs(slog.String("request_id", id))
	}
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		out.AddAttrs(slog.String("trace_id", sc.TraceID().String()), slog.String("span_id", sc.SpanID().String()))
	}
	return h.next.Handle(ctx, out)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clean := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		clean[i] = scrub(a)
	}
	return &handler{next: h.next.WithAttrs(clean)}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{next: h.next.WithGroup(name)}
}

func scrub(a slog.Attr) slog.Attr {
	if IsSensitive(a.Key) {
		return slog.String(a.Key, Redacted)
	}
	v := a.Value.Resolve()
	if v.Kind() != slog.KindGroup {
		return slog.Attr{Key: a.Key, Value: v}
	}
	group := v.Group()
	clean := make([]any, len(group))
	for i, g := range group {
		clean[i] = scrub(g)
	}
	return slog.Group(a.Key, clean...)
}
