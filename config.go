package log

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type option uint8

const (
	WithArchive option = iota + 1
	WithOutArchive
)
const ErrConfig option = 1<<8 - 1

var separator = string(filepath.Separator)

type Config struct {
	zap.Config
	ArchConf *ArchiveConfig
	Stdout   bool
}

// ArchiveConfig 日志文件归档配置
type ArchiveConfig struct {
	Dir        string
	ErrorFile  string
	WarnFile   string
	InfoFile   string
	DebugFile  string
	infoLog    *lumberjack.Logger
	errLog     *lumberjack.Logger
	warnLog    *lumberjack.Logger
	debugLog   *lumberjack.Logger
	MaxSize    int  // 按大小切割（M）
	MaxBackups int  // 默认备份数
	MaxAge     int  // 保存的最大天数
	Compress   bool // 是否对日志进行压缩
}

func (c *Config) build() (*zap.Logger, error) {
	switch c.check() {
	case WithArchive:
		c.buildArchLogger()
		return c.Config.Build(c.wrapSyncOpts())
	case WithOutArchive:
		return c.Config.Build()
	default:
		return nil, errors.New("log config error")
	}
}

func (c *Config) check() option {
	if c.ArchConf == nil || reflect.DeepEqual(*c.ArchConf, ArchiveConfig{}) {
		return WithOutArchive
	}
	if c.ArchConf.Dir == "" {
		c.ArchConf.Dir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
		c.ArchConf.Dir += separator + "logs" + separator
	}
	c.ArchConf.Dir = strings.TrimSuffix(c.ArchConf.Dir, separator)
	if c.ArchConf.DebugFile != "" ||
		c.ArchConf.InfoFile != "" ||
		c.ArchConf.WarnFile != "" ||
		c.ArchConf.ErrorFile != "" {
		c.ArchConf.DebugFile = strings.TrimSuffix(strings.Trim(c.ArchConf.DebugFile, separator), separator)
		c.ArchConf.InfoFile = strings.TrimSuffix(strings.Trim(c.ArchConf.InfoFile, separator), separator)
		c.ArchConf.WarnFile = strings.TrimSuffix(strings.Trim(c.ArchConf.WarnFile, separator), separator)
		c.ArchConf.ErrorFile = strings.TrimSuffix(strings.Trim(c.ArchConf.ErrorFile, separator), separator)
		return WithArchive
	}
	return ErrConfig
}

func (c *Config) buildArchLogger() {
	if c.ArchConf.InfoFile != "" {
		c.ArchConf.buildInfoLog()
	}
	if c.ArchConf.DebugFile != "" {
		c.ArchConf.buildDebugLog()
	}
	if c.ArchConf.ErrorFile != "" {
		c.ArchConf.buildErrorLog()
	}
	if c.ArchConf.WarnFile != "" {
		c.ArchConf.buildWarnLog()
	}
}

func (c *Config) wrapSyncOpts() zap.Option {
	fileEncoder := zapcore.NewJSONEncoder(c.Config.EncoderConfig)
	cores := make([]zapcore.Core, 0, 4)
	errEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel && zapcore.ErrorLevel-c.Config.Level.Level() > -1
	})
	warnEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-c.Config.Level.Level() > -1
	})
	infoEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-c.Config.Level.Level() > -1
	})
	debugEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-c.Config.Level.Level() > -1
	})
	if c.ArchConf.errLog != nil {
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(c.ArchConf.errLog), errEnabler))
	}

	if c.ArchConf.warnLog != nil {
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(c.ArchConf.warnLog), warnEnabler))
	}

	if c.ArchConf.infoLog != nil {
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(c.ArchConf.infoLog), infoEnabler))
	}

	if c.ArchConf.debugLog != nil {
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(c.ArchConf.debugLog), debugEnabler))
	}

	encoder := c.Config.EncoderConfig
	encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoder)
	if c.Stdout {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), errEnabler),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), warnEnabler),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), infoEnabler),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), debugEnabler),
		}...)
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}

func (c *Config) SetLevel(level string) {
	var lvl zapcore.Level
	switch level {
	case "debug", "Debug", "DEBUG":
		lvl = zapcore.DebugLevel
	case "info", "Info", "INFO":
		lvl = zapcore.InfoLevel
	case "warn", "Warn", "WARN":
		lvl = zapcore.WarnLevel
	case "error", "Error", "ERROR":
		lvl = zapcore.ErrorLevel
	case "panic", "Panic", "PANIC":
		lvl = zapcore.PanicLevel
	case "fatal", "Fatal", "FATAL":
		lvl = zapcore.FatalLevel
	default:
		lvl = zapcore.InfoLevel
	}
	c.Config.Level.SetLevel(lvl)
}

func (arch *ArchiveConfig) buildWarnLog() {
	arch.warnLog = &lumberjack.Logger{
		Filename:   arch.Dir + string(filepath.Separator) + arch.WarnFile,
		MaxSize:    arch.MaxSize,
		MaxBackups: arch.MaxBackups,
		LocalTime:  true,
		Compress:   arch.Compress,
	}
}

func (arch *ArchiveConfig) buildInfoLog() {
	arch.infoLog = &lumberjack.Logger{
		Filename:   arch.Dir + string(filepath.Separator) + arch.InfoFile,
		MaxSize:    arch.MaxSize,
		MaxBackups: arch.MaxBackups,
		LocalTime:  true,
		Compress:   arch.Compress,
	}
}

func (arch *ArchiveConfig) buildDebugLog() {
	arch.debugLog = &lumberjack.Logger{
		Filename:   arch.Dir + string(filepath.Separator) + arch.DebugFile,
		MaxSize:    arch.MaxSize,
		MaxBackups: arch.MaxBackups,
		LocalTime:  true,
		Compress:   arch.Compress,
	}
}

func (arch *ArchiveConfig) buildErrorLog() {
	arch.errLog = &lumberjack.Logger{
		Filename:   arch.Dir + string(filepath.Separator) + arch.ErrorFile,
		MaxSize:    arch.MaxSize,
		MaxBackups: arch.MaxBackups,
		LocalTime:  true,
		Compress:   arch.Compress,
	}
}

func NewDefaultConfig() Config {
	zapConf := zap.NewDevelopmentConfig()
	zapConf.EncoderConfig = DefaultZapEncoderConfig()
	archConf := DefaultArchiveConfig()
	return Config{
		Config:   zapConf,
		ArchConf: archConf,
		Stdout:   true,
	}
}

func NewDefaultWithoutArchiveConfig() Config {
	zapConf := zap.NewDevelopmentConfig()
	zapConf.EncoderConfig = DefaultZapEncoderConfig()
	return Config{
		Config: zapConf,
		Stdout: true,
	}
}

func DefaultZapEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "content",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func DefaultArchiveConfig() *ArchiveConfig {
	return &ArchiveConfig{
		Dir:        "",
		ErrorFile:  "error.log",
		WarnFile:   "warn.log",
		InfoFile:   "info.log",
		DebugFile:  "debug.log",
		MaxSize:    11,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   false,
	}
}
