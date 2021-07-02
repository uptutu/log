package log

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func New() *zap.Logger {
	config := NewDefaultConfig()
	var err error
	logger, err = config.build()
	if err != nil {
		panic(err)
	}
	return logger
}

func NewLogger(conf Config) *zap.Logger {
	var err error
	logger, err = conf.build()
	if err != nil {
		panic(err)
	}
	return logger
}

func Get() *zap.Logger {
	return logger
}

func Set(log *zap.Logger) {
	logger = log
}
