package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/library/errs"
	"github.com/labstack/echo/v4"
)

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			ctx := req.Context()

			var reqBody = make(map[string]any)

			reqContentType := req.Header.Get(echo.HeaderContentType)
			if strings.Contains(reqContentType, echo.MIMEApplicationJSON) && req.Body != nil {
				rawReqBody, err := io.ReadAll(req.Body)
				if err != nil {
					slog.ErrorContext(ctx, err.Error())
					resp := errs.ErrAPIFailed
					return c.JSON(resp.Status, resp)
				}
				// ---------------------------
				// TODO mask with regexp here
				// ---------------------------
				if err := json.Unmarshal(rawReqBody, &reqBody); err != nil {
					slog.LogAttrs(ctx, slog.LevelError, err.Error())
					resp := errs.ErrAPIFailed
					return c.JSON(resp.Status, resp)
				}
				// Set the body back to the request.
				c.Request().Body = io.NopCloser(bytes.NewBuffer(rawReqBody))
			}

			slog.LogAttrs(ctx, slog.LevelInfo, fmt.Sprintf("request %s %s", req.Method, req.URL.String()),
				slog.String("contentType", reqContentType),
				slog.Any("reqBody", reqBody),
				slog.String("host", req.Host),
				slog.String("ip", c.RealIP()),
				slog.Any("formValue", req.Form),
			)

			// Response
			rawResBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, rawResBody)
			writer := &MsgResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			res := c.Response()
			resContentType := res.Header().Get(echo.HeaderContentType)

			var resBody = make(map[string]any)
			if strings.Contains(resContentType, echo.MIMEApplicationJSON) && rawResBody != nil {
				if err := json.Unmarshal(rawResBody.Bytes(), &resBody); err != nil {
					slog.LogAttrs(ctx, slog.LevelError, err.Error())
					resp := errs.ErrAPIFailed
					return c.JSON(resp.Status, resp)
				}
			}

			var level slog.Level
			switch {
			case res.Status >= 500:
				level = slog.LevelError
			case res.Status >= 400:
				level = slog.LevelWarn
			case res.Status >= 300:
				level = slog.LevelInfo
			default:
				level = slog.LevelInfo
			}

			message := fmt.Sprintf("response %s %s", req.Method, req.URL.String())

			if strings.Contains(resContentType, "HTML") {
				slog.LogAttrs(ctx, level, message,
					slog.Int("code", res.Status),
					slog.String("contentType", resContentType),
					slog.Duration("latency", time.Since(start)),
					slog.Any("data", c.Get(string(ContextKeyHTMXData))),
				)
			} else {
				slog.LogAttrs(ctx, level, message,
					slog.Int("code", res.Status),
					slog.String("contentType", resContentType),
					slog.Any("resBody", resBody),
					slog.Duration("latency", time.Since(start)),
				)
			}

			return nil
		}
	}
}

type MsgResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *MsgResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *MsgResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
