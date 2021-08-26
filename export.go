package log

import "go.uber.org/zap"

// export Debug, Info, Warn, Error, Panic, Fatal shortcut

// Debug will log at debug level
func Debug(msg string, fields ...zap.Field) {
	callerSkippedLog().Debug(msg, fields...)
}

// Info logs a message at InfoLevel
func Info(msg string, fields ...zap.Field) {
	callerSkippedLog().Info(msg, fields...)
}

// Warn will log stacktrace info
func Warn(msg string, fields ...zap.Field) {
	callerSkippedLog().Warn(msg, fields...)
}

// Error will log stacktrace info
func Error(msg string, fields ...zap.Field) {
	callerSkippedLog().Error(msg, fields...)
}

// Panic The logger then panics, even if logging at PanicLevel is disabled, will recovery if set
func Panic(msg string, fields ...zap.Field) {
	callerSkippedLog().Panic(msg, fields...)
}

// Fatal The logger then calls os.Exit(1)
func Fatal(msg string, fields ...zap.Field) {
	callerSkippedLog().Fatal(msg, fields...)
}

func callerSkippedLog() *zap.Logger {
	return zap.L().WithOptions(zap.AddCallerSkip(1))
}
