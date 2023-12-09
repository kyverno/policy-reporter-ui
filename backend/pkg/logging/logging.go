package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Enabled     bool   `mapstructure:"enabled"`
	LogLevel    int8   `mapstructure:"logLevel"`
	Encoding    string `mapstructure:"encoding"`
	Development bool   `mapstructure:"development"`
}

func New(config Config) *zap.Logger {
	encoder := zap.NewProductionEncoderConfig()
	if config.Development {
		encoder = zap.NewDevelopmentEncoderConfig()
		encoder.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	}

	ouput := "json"
	if config.Encoding != "json" {
		ouput = "console"
	}

	var sampling *zap.SamplingConfig
	if !config.Development {
		sampling = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}
	}

	cnfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.Level(config.LogLevel)),
		Development:       config.Development,
		Sampling:          sampling,
		Encoding:          ouput,
		EncoderConfig:     encoder,
		DisableStacktrace: !config.Development,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	logger, _ := cnfg.Build()

	zap.ReplaceGlobals(logger)

	return logger
}
