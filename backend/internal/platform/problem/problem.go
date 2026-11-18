// Package problem builds and writes RFC 9457 problem details (application/problem+json).
package problem

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/an4eetos/inner_circle/backend/internal/api"
	"github.com/an4eetos/inner_circle/backend/internal/platform/reqctx"
)

// ContentType is the media type for problem responses.
const ContentType = "application/problem+json"

// Stable error codes. Clients switch on these, never on titles or details.
const (
	CodeBadRequest       = "bad_request"
	CodeNotFound         = "not_found"
	CodePayloadTooLarge  = "payload_too_large"
	CodeRateLimited      = "rate_limited"
	CodeInternal         = "internal"
	CodeUnavailable      = "unavailable"
	CodeMethodNotAllowed = "method_not_allowed"
)

// New builds a problem for status with a stable code; detail is optional ("" omits it).
// The request ID from ctx is attached so clients can quote it in bug reports.
func New(ctx context.Context, status int, code, detail string) api.Problem {
	p := api.Problem{
		Type:   "about:blank",
		Title:  http.StatusText(status),
		Status: status,
		Code:   code,
	}
	if detail != "" {
		p.Detail = &detail
	}
	if id := reqctx.RequestID(ctx); id != "" {
		p.RequestId = &id
	}
	return p
}

// Write sends a problem response.
func Write(w http.ResponseWriter, r *http.Request, status int, code, detail string) {
	p := New(r.Context(), status, code, detail)
	w.Header().Set("Content-Type", ContentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(p)
}
