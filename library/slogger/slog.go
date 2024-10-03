package slogger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bosskrub9992/fuel-management-backend/library/masks"
)

const LevelCritical = slog.Level(12)

type Config struct {
	IsProductionEnv bool
	MaskingFields   []string
	RemovingFields  []string
}

func New(cfg *Config) *slog.Logger {
	minimumLogLevel := slog.LevelDebug
	if cfg.IsProductionEnv {
		minimumLogLevel = slog.LevelInfo
	}

	var maskingField = make(map[string]bool)
	for _, f := range cfg.MaskingFields {
		maskingField[f] = true
	}

	var removingField = make(map[string]bool)
	for _, f := range cfg.RemovingFields {
		removingField[f] = true
	}

	replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
		if _, found := maskingField[a.Key]; found {
			if a.Value.Kind() == slog.KindString {
				a.Value = slog.StringValue(masks.Left(a.Value.String(), 4))
				return a
			}
		}
		if _, found := removingField[a.Key]; found {
			return slog.Attr{}
		}
		if !cfg.IsProductionEnv {
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				a.Value = slog.StringValue(fmt.Sprintf("%s:%d", source.File, source.Line))
				return a
			}
		}
		if a.Key == slog.MessageKey {
			a.Key = "message"
		} else if a.Key == slog.SourceKey {
			a.Key = "logging.googleapis.com/sourceLocation"
		} else if a.Key == slog.LevelKey {
			a.Key = "severity"
			level := a.Value.Any().(slog.Level)
			if level == LevelCritical {
				a.Value = slog.StringValue("CRITICAL")
			}
		}
		return a
	}

	var slogHandler slog.Handler
	if cfg.IsProductionEnv {
		slogHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,
			Level:       minimumLogLevel,
			ReplaceAttr: replaceAttr,
		})
	} else {
		slogHandler = newPrettyHandler(&slog.HandlerOptions{
			AddSource:   true,
			Level:       minimumLogLevel,
			ReplaceAttr: replaceAttr,
		})
	}
	slogHandler = newContextHandler(slogHandler)

	logger := slog.New(slogHandler)
	slog.SetDefault(logger)
	return logger
}
