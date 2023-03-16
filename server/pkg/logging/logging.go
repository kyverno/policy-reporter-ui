package logging

import (
	"github.com/kyverno/policy-reporter-ui/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(config *config.Config) *zap.Logger {
	encoder := zap.NewProductionEncoderConfig()
	if config.Logging.Development {
		encoder = zap.NewDevelopmentEncoderConfig()
		encoder.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	}

	ouput := "json"
	if config.Logging.Encoding != "json" {
		ouput = "console"
	}

	var sampling *zap.SamplingConfig
	if !config.Logging.Development {
		sampling = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}
	}

	cnfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.Level(config.Logging.LogLevel)),
		Development:       config.Logging.Development,
		Sampling:          sampling,
		Encoding:          ouput,
		EncoderConfig:     encoder,
		DisableStacktrace: true,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	logger, _ := cnfg.Build()

	return logger
}
