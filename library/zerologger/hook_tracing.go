package zerologger

import (
	"context"

	"github.com/bosskrub9992/fuel-management-backend/library/middlewares"
	"github.com/rs/zerolog"
)

type TracingHook struct{}

func (h TracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	traceID := getRequestIDFromContext(ctx)
	e.Str(string(middlewares.LogFieldKeyReqID), traceID)
}

func getRequestIDFromContext(ctx context.Context) string {
	reqID, ok := ctx.Value(middlewares.ContextKeyRequestID).(string)
	if !ok {
		return ""
	}
	return reqID
}
