package slogger

import (
	"context"
	"log/slog"

	"github.com/bosskrub9992/fuel-management-backend/library/middlewares"
)

type ContextHandler struct {
	h slog.Handler
}

func newContextHandler(handler slog.Handler) *ContextHandler {
	return &ContextHandler{
		h: handler,
	}
}

func (h ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(h.addTraceFromContext(ctx)...)
	return h.h.Handle(ctx, r)
}

func (h ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.h.WithAttrs(attrs)
}

func (h ContextHandler) WithGroup(name string) slog.Handler {
	return h.h.WithGroup(name)
}

func (h ContextHandler) addTraceFromContext(ctx context.Context) []slog.Attr {
	attrs := []slog.Attr{}
	reqID, ok := ctx.Value(middlewares.ContextKeyRequestID).(string)
	if ok {
		attrs = append(attrs, slog.String(
			string(middlewares.LogFieldKeyReqID), reqID),
		)
	}
	return attrs
}
