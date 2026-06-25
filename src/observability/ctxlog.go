// Package observability provides correlation-aware structured logging for the chat
// pipeline: a slog handler that injects request/org/user ids taken from the context
// into every record, plus a verbosity toggle that gates full-text (PII) logging.
//
// It is a leaf package (only stdlib + uuid) so core usecases and infra repos can
// import it without breaking clean-architecture boundaries.
package observability

import (
	"context"
	"log/slog"
	"sync/atomic"

	"github.com/google/uuid"
)

type ctxKey int

const (
	requestIDKey ctxKey = iota
	orgIDKey
	userIDKey
)

// WithRequestID/WithOrg/WithUser stash correlation ids in the context so the
// ContextHandler can attach them to every log emitted with slog.*Context(ctx, ...).
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func WithOrg(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, orgIDKey, id)
}

func WithUser(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// ContextHandler wraps a slog.Handler and enriches each record with the correlation
// ids present in the context (request_id, org_id, user_id). Missing ids are skipped.
type ContextHandler struct {
	slog.Handler
}

func NewContextHandler(h slog.Handler) *ContextHandler {
	return &ContextHandler{Handler: h}
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if v, ok := ctx.Value(requestIDKey).(string); ok && v != "" {
		r.AddAttrs(slog.String("request_id", v))
	}
	if v, ok := ctx.Value(orgIDKey).(uuid.UUID); ok && v != uuid.Nil {
		r.AddAttrs(slog.String("org_id", v.String()))
	}
	if v, ok := ctx.Value(userIDKey).(int64); ok && v != 0 {
		r.AddAttrs(slog.Int64("user_id", v))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{Handler: h.Handler.WithGroup(name)}
}

// verbose gates full-text (PII) logging. Default true; set from config at boot.
var verbose atomic.Bool

func init() { verbose.Store(true) }

// SetVerbose enables/disables full-text logging (e.g. prompts, responses, names).
func SetVerbose(v bool) { verbose.Store(v) }

// Verbose reports whether full-text logging is on.
func Verbose() bool { return verbose.Load() }

// Text returns the full value as attr `key` when verbose; otherwise it returns
// `key+"_len"` with the rune length, so the same log call carries full text in dev
// and only a size in prod without leaking PII.
func Text(key, value string) slog.Attr {
	if verbose.Load() {
		return slog.String(key, value)
	}
	return slog.Int(key+"_len", len([]rune(value)))
}
