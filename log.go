package log

import (
	"go.uber.org/zap"
	"reflect"
)

func init() {
	if reflect.DeepEqual(zap.L(), zap.NewNop()) {
		Set(NewConsoleLogger())
	}
}

func New() *zap.Logger {
	config := NewDefaultConfig()
	logger, err := config.build()
	if err != nil {
		panic(err)
	}
	return logger
}

func NewConsoleLogger() *zap.Logger {
	config := NewDefaultWithoutArchiveConfig()
	logger, err := config.build()
	if err != nil {
		panic(err)
	}
	return logger
}

func NewLogger(conf Config) *zap.Logger {
	logger, err := conf.build()
	if err != nil {
		panic(err)
	}
	return logger
}

func Get() *zap.Logger {
	return zap.L()
}

func Set(log *zap.Logger) {
	zap.ReplaceGlobals(log)
}

func SetLogWithFields(m map[string]string) {
	Set(WrapFields(m))
}

func WrapFields(m map[string]string) *zap.Logger {
	fields := make([]zap.Option, 0, len(m))
	for k, v := range m {
		fields = append(fields, zap.Fields(zap.String(k, v)))
	}

	return zap.L().WithOptions(fields...)
}
