package cmd

import (
	"context"
	"errors"
	"grubzo/internal/utils"
	"time"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLogger() (logger *zap.Logger, cleanup func(context.Context) error) {
	if c.DevMode {
		return getCLILogger()
	}
	// cfg := zap.Config{
	// 	Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
	// 	Encoding:         "json",
	// 	EncoderConfig:    zapdriver.NewProductionEncoderConfig(),
	// 	OutputPaths:      []string{"stdout"},
	// 	ErrorOutputPaths: []string{"stderr"},
	// }
	// logger, _ = cfg.Build(zapdriver.WrapCore(zapdriver.ServiceName(fmt.Sprintf("grubzo.%s.%s", Version, Revision))))

	lokiWriter := utils.NewLokiWriter(
		c.LokiHost,
		map[string]string{
			"job":     "GrubzoServer",
			"service": "grubzo",
			"env":     "prod",
		},
		2*time.Second,
	)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapdriver.NewProductionEncoderConfig()),
		zapcore.AddSync(lokiWriter),
		zap.InfoLevel,
	)
	logger = zap.New(core)
	cleanup = func(ctx context.Context) error {
		return errors.Join(
			logger.Sync(),
			lokiWriter.Close(ctx),
		)
	}
	return logger, cleanup
}

func getCLILogger() (logger *zap.Logger, cleanup func(context.Context) error) {
	level := zap.NewAtomicLevel()
	if c.DevMode {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	cfg := zap.Config{
		Level:       level,
		Development: c.DevMode,
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
