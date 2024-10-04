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

func init() {
	zerolog.LevelFieldName = "severity"
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		// mapping to Cloud Logging LogSeverity
		// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
		// reference: https://github.com/yfuruyama/crzerolog/blob/master/log.go
		switch l {
		case zerolog.TraceLevel:
			return "DEFAULT"
		case zerolog.DebugLevel:
			return "DEBUG"
		case zerolog.InfoLevel:
			return "INFO"
		case zerolog.WarnLevel:
			return "WARNING"
		case zerolog.ErrorLevel:
			return "ERROR"
		case zerolog.FatalLevel:
			return "CRITICAL"
		case zerolog.PanicLevel:
			return "ALERT"
		case zerolog.NoLevel:
			return "DEFAULT"
		default:
			return "DEFAULT"
		}
	}
	zerolog.CallerFieldName = "logging.googleapis.com/sourceLocation"
}
