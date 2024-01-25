package main

import (
	"log/slog"

	"github.com/bosskrub9992/fuel-management-backend/config"
	"github.com/bosskrub9992/fuel-management-backend/library/slogger"
)

func main() {
	cfg := config.New()
	slog.SetDefault(slogger.New(&slogger.Config{
		IsProductionEnv: false,
	}))

	slog.Info("test config", slog.Any("cfg", cfg))
}
