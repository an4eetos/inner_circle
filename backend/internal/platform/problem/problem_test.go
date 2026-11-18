package problem

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/an4eetos/inner_circle/backend/internal/api"
	"github.com/an4eetos/inner_circle/backend/internal/platform/reqctx"
)

func TestWrite(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/x", nil)
	r = r.WithContext(reqctx.WithRequestID(r.Context(), "req-42"))
	w := httptest.NewRecorder()

	Write(w, r, http.StatusNotFound, CodeNotFound, "no such thing")

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != ContentType {
		t.Fatalf("content-type = %q", ct)
	}
	var p api.Problem
	if err := json.Unmarshal(w.Body.Bytes(), &p); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if p.Code != CodeNotFound || p.Title != "Not Found" || p.Status != 404 || p.Type != "about:blank" ||
		p.Detail == nil || *p.Detail != "no such thing" || p.RequestId == nil || *p.RequestId != "req-42" {
		t.Fatalf("unexpected problem: %+v", p)
	}
}

func TestNew_OmitsEmptyDetail(t *testing.T) {
	p := New(httptest.NewRequest(http.MethodGet, "/", nil).Context(), http.StatusInternalServerError, CodeInternal, "")
	if p.Detail != nil || p.RequestId != nil {
		t.Fatalf("expected nil detail and request id: %+v", p)
	}
}
