package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/signal"
	"text/template"
	"time"

	"log/slog"

	"github.com/bosskrub9992/fuel-management/config"
	"github.com/bosskrub9992/fuel-management/internal/adaptors/gormadaptor"
	"github.com/bosskrub9992/fuel-management/internal/handlers/htmxhandler"
	"github.com/bosskrub9992/fuel-management/internal/handlers/resthandler"
	"github.com/bosskrub9992/fuel-management/internal/services"
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
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	dataInBytes, err := json.Marshal(data)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}
	c.Set("data", string(dataInBytes))

	return nil
}

func main() {
	// dependency injection
	cfg := config.New()
	logger := slogger.New(&cfg.Logger)
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
	restHandler := resthandler.NewRESTHandler(healthService)
	htmxHandler := htmxhandler.NewHTMXHandler(service)

	e := echo.New()
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("internal/templates/src/*.html")),
	}
	e.Static("/dist", "./internal/templates/dist")
	e.Static("/node_modules", "./internal/templates/node_modules")
	e.Static("/static", "./internal/templates/static")
	e.Use(
		middleware.Recover(),
		middleware.CORS(),
		slogger.MiddlewareHTMX(),
	)
	apiV1Group := e.Group("/api/v1", slogger.MiddlewareREST())
	apiV1Group.GET("/health", restHandler.GetHealth, slogger.MiddlewareREST())

	e.GET("/fuel-usage", htmxHandler.FuelUsage)
	e.POST("/create-fuel-usage", htmxHandler.CreateFuelUsage)
	e.GET("/test", htmxHandler.Test)
	e.POST("/example", htmxHandler.Example)

	// run server
	go func() {
		if err := e.Start(":" + cfg.Server.Port); err != nil && err != http.ErrServerClosed {
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
