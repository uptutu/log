package log

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	zap.Config
	ArchConf ArchiveConfig
	Stdout   bool // 是否在控制台输出
}

// ArchiveConfig 日志文件归档配置
type ArchiveConfig struct {
	LogFileDir    string //文件保存目录
	ErrorFileName string
	WarnFileName  string
	InfoFileName  string
	DebugFileName string
	infoLog       *lumberjack.Logger
	errLog        *lumberjack.Logger
	warnLog       *lumberjack.Logger
	debugLog      *lumberjack.Logger
	MaxSize       int  // 按大小切割（M）
	MaxBackups    int  // 默认备份数
	MaxAge        int  // 保存的最大天数
	Compress      bool // 是否对日志进行压缩
}

func (c *Config) build() (*zap.Logger, error) {
	if !c.check() {
		return nil, errors.New("config error")
	}
	c.buildArchLogger()
	opts := c.wrapSync()

	return c.Config.Build(opts)
}

func (c Config) check() bool {
	if c.ArchConf.LogFileDir == "" {
		c.ArchConf.LogFileDir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
		c.ArchConf.LogFileDir += string(filepath.Separator) + "logs" + string(filepath.Separator)
	}
	c.ArchConf.LogFileDir = strings.TrimSuffix(c.ArchConf.LogFileDir, string(filepath.Separator))
	return true
}

func (c *Config) buildArchLogger() {
	c.ArchConf.errLog = c.ArchConf.getLumberjackLogger(c.ArchConf.ErrorFileName)
	c.ArchConf.infoLog = c.ArchConf.getLumberjackLogger(c.ArchConf.InfoFileName)
	c.ArchConf.debugLog = c.ArchConf.getLumberjackLogger(c.ArchConf.DebugFileName)
	c.ArchConf.warnLog = c.ArchConf.getLumberjackLogger(c.ArchConf.WarnFileName)
}

func (c *Config) wrapSync() zap.Option {
	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel && zapcore.ErrorLevel-c.Config.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-c.Config.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-c.Config.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-c.Config.Level.Level() > -1
	})

	fileEncoder := zapcore.NewJSONEncoder(c.Config.EncoderConfig)
	cores := []zapcore.Core{
		zapcore.NewCore(fileEncoder, zapcore.AddSync(c.ArchConf.errLog), errPriority),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(c.ArchConf.warnLog), warnPriority),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(c.ArchConf.infoLog), infoPriority),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(c.ArchConf.debugLog), debugPriority),
	}

	c.Config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(c.Config.EncoderConfig)
	if c.Stdout {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), errPriority),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), warnPriority),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), infoPriority),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), debugPriority),
		}...)
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}

func (arch ArchiveConfig) getLumberjackLogger(name string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   arch.LogFileDir + string(filepath.Separator) + name,
		MaxSize:    arch.MaxSize,
		MaxAge:     arch.MaxAge,
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
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func (arch *ArchiveConfig) StopArchive(kindName string) error {
	switch strings.ToLower(kindName) {
	case "info", "information":
		return arch.infoLog.Close()
	case "err", "error", "errors":
		return arch.errLog.Close()
	case "debug":
		return arch.debugLog.Close()
	case "warn", "warning":
		return arch.warnLog.Close()
	}

	return errors.New("invalid kind close")
}

func (arch *ArchiveConfig) Rotate(kindName string) error {
	switch strings.ToLower(kindName) {
	case "info", "information":
		return arch.infoLog.Rotate()
	case "err", "error", "errors":
		return arch.errLog.Rotate()
	case "debug":
		return arch.debugLog.Rotate()
	case "warn", "warning":
		return arch.warnLog.Rotate()
	}

	return errors.New("invalid kind rotate")
}

func DefaultArchiveConfig() ArchiveConfig {
	return ArchiveConfig{
		LogFileDir:    "",
		ErrorFileName: "error.log",
		WarnFileName:  "warn.log",
		InfoFileName:  "info.log",
		DebugFileName: "debug.log",
		MaxSize:       11,
		MaxBackups:    10,
		MaxAge:        30,
		Compress:      false,
	}
}
