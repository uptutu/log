package log

import (
	"go.uber.org/zap"
)

func New() *zap.Logger {
	config := NewDefaultConfig()
	logger, err := config.build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	return logger
}

func NewLogger(conf Config) *zap.Logger {
	logger, err := conf.build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	return logger
}

func Get() *zap.Logger {
	return zap.L()
}

func Set(log *zap.Logger) {
	zap.ReplaceGlobals(log)
}
