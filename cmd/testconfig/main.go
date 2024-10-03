package main

import (
	"github.com/bosskrub9992/fuel-management-backend/config"
	"github.com/bosskrub9992/fuel-management-backend/library/zerologger"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.New()
	zerologger.InitZerologExtension(cfg.Logger)

	log.Info().Any("cfg", cfg).Msg("test config")
}
