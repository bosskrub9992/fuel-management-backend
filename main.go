package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"text/template"
	"time"

	"github.com/bosskrub9992/fuel-management/config"
	"github.com/bosskrub9992/fuel-management/internal/adaptors/gormadaptor"
	"github.com/bosskrub9992/fuel-management/internal/handlers/resthandler"
	"github.com/bosskrub9992/fuel-management/internal/services"

	"log/slog"

	"github.com/jinleejun-corp/corelib/databases"
	"github.com/jinleejun-corp/corelib/slogger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	ctx := c.Request().Context()
	if err := t.templates.ExecuteTemplate(w, name, data); err != nil {
		slogger.Error(ctx, err.Error())
		return err
	}
	return nil
}

func main() {
	// dependency injection
	cfg := config.New()
	logger := slogger.Init(&cfg.Logger)
	slog.SetDefault(logger)
	sqlDB, err := databases.NewPostgres(&cfg.Database)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer sqlDB.Close()
	gormDB, err := databases.NewGormDBPostgres(sqlDB)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	db := gormadaptor.NewDatabase(gormDB)
	healthService := services.NewHealthService()
	service := services.New(cfg, db)
	handler := resthandler.NewRESTHandler(healthService, service)

	e := echo.New()
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("internal/templates/src/*.html")),
	}
	e.Static("/dist", "./internal/templates/dist")
	e.Static("/node_modules", "./internal/templates/node_modules")
	e.Use(
		middleware.Recover(),
		middleware.CORS(),
		slogger.EchoMiddleware(),
	)
	e.GET("/fuel-usage", handler.FuelUsage) // TODO update middleware to support form value
	apiV1 := e.Group("/api/v1")
	apiV1.GET("/health", handler.GetHealth)
	apiV1.POST("/customers", handler.CreateCustomer)

	// run rest server
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		panic(err)
	}
}
