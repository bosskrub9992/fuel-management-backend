package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/config"
	"github.com/bosskrub9992/fuel-management-backend/internal/adaptors/pgadaptor"
	"github.com/bosskrub9992/fuel-management-backend/internal/handlers/resthandler"
	"github.com/bosskrub9992/fuel-management-backend/internal/routers"
	"github.com/bosskrub9992/fuel-management-backend/internal/services"
	"github.com/jinleejun-corp/corelib/databases"
	"github.com/jinleejun-corp/corelib/slogger"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.New()
	slog.SetDefault(slogger.New(&cfg.Logger))
	sqlDB, err := databases.NewPostgres(&cfg.Database.Postgres)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()
	gormDB, err := databases.NewGormDBPostgres(sqlDB)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	db := pgadaptor.NewPostgresAdaptor(gormDB)
	service := services.New(cfg, db)
	restHandler := resthandler.New(service)

	e := echo.New()
	router := routers.New(e, restHandler)
	e = router.Init()

	// run server
	go func() {
		address := fmt.Sprintf(":%s", cfg.Server.Port)
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			slog.Error(err.Error())
			return
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		slog.Error(err.Error())
		return
	}
}
