package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLogger() (*zap.Logger, func(context.Context) error) {
	if c.IsDev() {
		return getCLILogger()
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding:         "json",
		EncoderConfig:    zapdriver.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build(
		zapdriver.WrapCore(
			zapdriver.ServiceName(
				fmt.Sprintf("grubzo.%s.%s", Version, Revision),
			),
		),
	)
	if err != nil {
		panic(err)
	}

	cleanup := func(ctx context.Context) error {
		return errors.Join(
			logger.Sync(),
		)
	}

	return logger, cleanup
}

func getCLILogger() (logger *zap.Logger, cleanup func(context.Context) error) {
	level := zap.NewAtomicLevel()
	if c.IsDev() {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	cfg := zap.Config{
		Level:       level,
		Development: c.IsDev(),
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ = cfg.Build()
	cleanup = func(context.Context) error {
		return logger.Sync()
	}
	return logger, cleanup
}
