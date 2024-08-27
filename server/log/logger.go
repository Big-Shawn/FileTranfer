package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"server/config"
)

var L *zap.Logger

func Init() error {
	c := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.Level(config.C.Logger.Level)),
		Development:      config.C.Logger.Development,
		Encoding:         config.C.Logger.Encoding,
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      config.C.Logger.Path,
		ErrorOutputPaths: config.C.Logger.Path,
	}
	logger, err := c.Build()
	L = logger
	if err != nil {
		return err
	}

	return nil
}
