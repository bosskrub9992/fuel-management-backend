package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/library/errs"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ZeroLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			ctx := req.Context()

			var reqBody = []byte("{}")
			var err error

			reqContentType := req.Header.Get(echo.HeaderContentType)
			if strings.Contains(reqContentType, echo.MIMEApplicationJSON) && req.Body != nil {
				reqBody, err = io.ReadAll(req.Body)
				if err != nil {
					slog.ErrorContext(ctx, err.Error())
					resp := errs.ErrAPIFailed
					return c.JSON(resp.Status, resp)
				}

				// Set the body back to the request.
				c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}

			log.Info().Ctx(ctx).
				Str("contentType", reqContentType).
				RawJSON("reqBody", reqBody).
				Str("host", req.Host).
				Str("ip", c.RealIP()).
				Any("formValue", req.Form).
				Msgf("request %s %s", req.Method, req.URL.String())

			// Response
			rawResBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, rawResBody)
			writer := &MsgResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			err = next(c)
			if err != nil {
				c.Error(err)
			}

			res := c.Response()
			resContentType := res.Header().Get(echo.HeaderContentType)
			resBody := rawResBody.Bytes()

			var level zerolog.Level
			switch {
			case res.Status >= 500:
				level = zerolog.ErrorLevel
			case res.Status >= 400:
				level = zerolog.WarnLevel
			case res.Status >= 300:
				level = zerolog.InfoLevel
			default:
				level = zerolog.InfoLevel
			}

			message := fmt.Sprintf("response %s %s", req.Method, req.URL.String())

			if strings.Contains(resContentType, "HTML") {
				log.WithLevel(level).Ctx(ctx).
					Int("code", res.Status).
					Str("contentType", resContentType).
					Str("latency", time.Since(start).String()).
					Any("data", c.Get(string(ContextKeyHTMXData))).
					Msg(message)
			} else {
				log.WithLevel(level).Ctx(ctx).
					Int("code", res.Status).
					Str("contentType", resContentType).
					RawJSON("resBody", resBody).
					Str("latency", time.Since(start).String()).
					Msg(message)
			}

			return nil
		}
	}
}
