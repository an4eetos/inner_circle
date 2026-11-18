package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/an4eetos/inner_circle/backend/internal/platform/reqctx"
)

type secretValuer struct{}

func (secretValuer) LogValue() slog.Value {
	return slog.GroupValue(slog.String("medication", "Ursofalk 250mg"), slog.Int("count", 1))
}

func TestLogger_ScrubsSensitiveAttributes(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, slog.LevelDebug, "json").With("email", "a@b.test")

	ctx := reqctx.WithRequestID(context.Background(), "req-1")
	log.InfoContext(ctx, "check-in",
		"user_id", "u1",
		"Medication", "Ursofalk",
		"health_note", "gallstones",
		slog.Group("meal", "title", "fried chicken"),
		slog.Group("payload", "notes", "secret note", "kept", "visible"),
		"valuer", secretValuer{},
	)

	out := buf.String()
	for _, leaked := range []string{"a@b.test", "Ursofalk", "gallstones", "fried chicken", "secret note"} {
		if strings.Contains(out, leaked) {
			t.Fatalf("sensitive value %q leaked: %s", leaked, out)
		}
	}

	var rec map[string]any
	if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	if rec["user_id"] != "u1" || rec["request_id"] != "req-1" {
		t.Fatalf("expected user_id and request_id to survive: %v", rec)
	}
	payload, ok := rec["payload"].(map[string]any)
	if !ok || payload["kept"] != "visible" || payload["notes"] != Redacted {
		t.Fatalf("group scrubbing wrong: %v", rec["payload"])
	}
	if rec["email"] != Redacted || rec["meal"] != Redacted {
		t.Fatalf("expected redaction markers: %v", rec)
	}
}

func TestIsSensitive(t *testing.T) {
	for _, k := range []string{"TOKEN", "refresh_token", "health_anything", "Results"} {
		if !IsSensitive(k) {
			t.Errorf("%q should be sensitive", k)
		}
	}
	for _, k := range []string{"user_id", "status", "path", "duration_ms"} {
		if IsSensitive(k) {
			t.Errorf("%q should not be sensitive", k)
		}
	}
}
