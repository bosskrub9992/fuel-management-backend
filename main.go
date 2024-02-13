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
	"github.com/bosskrub9992/fuel-management-backend/library/databases"
	"github.com/bosskrub9992/fuel-management-backend/library/slogger"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func main() {
	cfg := config.New()
	slog.SetDefault(slogger.New(&slogger.Config{
		IsProductionEnv: cfg.Logger.IsProductionEnv,
		MaskingFields:   cfg.Logger.MaskingFields,
		RemovingFields:  cfg.Logger.RemovingFields,
	}))
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
	gormDB, err := databases.NewGormDBPostgres(sqlDB, gorm.Config{})
	if err != nil {
		slog.Error(err.Error())
		return
	}
	db, err := pgadaptor.NewPostgresAdaptor(gormDB)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	service, err := services.New(cfg, db)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	restHandler, err := resthandler.New(service, time.Now())
	if err != nil {
		slog.Error(err.Error())
		return
	}

	e := echo.New()
	router, err := routers.New(e, restHandler)
	if err != nil {
		slog.Error(err.Error())
		return
	}
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
