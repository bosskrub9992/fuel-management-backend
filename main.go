package main

import (
	"context"
	"fmt"
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
	"github.com/bosskrub9992/fuel-management-backend/library/zerologger"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func main() {
	cfg := config.New()
	zerologger.InitZerologExtension(cfg.Logger)
	sqlDB, err := databases.NewPostgres(&cfg.Database.Postgres)
	if err != nil {
		log.Info().Err(err).Send()
		return
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Info().Err(err).Send()
		}
	}()
	gormDB, err := databases.NewGormDBPostgres(sqlDB, gorm.Config{})
	if err != nil {
		log.Info().Err(err).Send()
		return
	}
	db := pgadaptor.NewPostgresAdaptor(gormDB)
	service := services.New(cfg, db)
	restHandler := resthandler.New(service, time.Now())

	e := echo.New()
	router := routers.New(e, restHandler)
	e = router.Init()

	// run server
	go func() {
		address := fmt.Sprintf(":%s", cfg.Server.Port)
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			log.Info().Err(err).Send()
			return
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Info().Msg("server is shuting down ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Info().Err(err).Send()
		return
	}
}
