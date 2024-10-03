package main

import (
	"net/http"

	"github.com/bosskrub9992/fuel-management-backend/config"
	"github.com/bosskrub9992/fuel-management-backend/library/middlewares"
	"github.com/bosskrub9992/fuel-management-backend/library/zerologger"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.New()

	zerologger.InitZerologExtension(cfg.Logger)

	e := echo.New()
	e.Use(
		middlewares.RequestID(),
		middlewares.ZeroLogger(),
	)
	e.GET("/test", handler)
	e.POST("/test", handler2)

	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Err(err).Msg("failed to close sever")
	}
}

func handler(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"hello": "world",
	})
}

func handler2(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"hello": "world",
	})
}
