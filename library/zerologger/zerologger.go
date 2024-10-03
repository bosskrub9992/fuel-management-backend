package zerologger

// TODO async logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ZerologConfig struct {
	IsProductionEnv bool `mapstructure:"is_production_env"`
}

func InitZerologExtension(cfg ZerologConfig) {
	if cfg.IsProductionEnv {
		log.Logger = log.Level(zerolog.InfoLevel)
	} else {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    false,
			TimeFormat: time.RFC3339,
		}).Level(zerolog.TraceLevel)
	}
	log.Logger = log.Logger.With().Caller().Timestamp().Logger().Hook(TracingHook{})
}
