package middlewares

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()

			reqID := req.Header.Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = uuid.NewString()
			}

			ctx = context.WithValue(ctx, ContextKeyRequestID, reqID)
			c.SetRequest(req.WithContext(ctx))

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			return nil
		}
	}
}
