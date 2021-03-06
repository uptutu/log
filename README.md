# Log

## Introduction

Wrap Uber-go /Zap make log easily.

## Installation

```shell
go get github.com/uptutu/log
```

## Usage

Easy way to use. This `New()` will set working directory as default config to this `logger` save log info.

```go
import (
"github.uptutu/log"
)

func main() {
log.New()
log.Info("xxxxx")
log.Debug("xxxxx")
}
```

Edit the config as you want, `Config` field is `uber-go/Zap` config, ArchConf as you want, if set `nil` to this field
then
`logger` disable archive.

```go
import (
"github.uptutu/log"
)

func main() {
conf := log.NewDefaultConfig()
conf.ArchConf = nil
log.NewLogger(conf)

log.Info("xxxxx")
log.Debug("xxxxx")
}
```

default config info:

```go
zap.Config{
    Level:            NewAtomicLevelAt(DebugLevel),
    Development:      true,
    Encoding:         "console",
    EncoderConfig:    NewDevelopmentEncoderConfig(),
    OutputPaths:      []string{"stderr"},
    ErrorOutputPaths: []string{"stderr"},
}
zapcore.EncoderConfig{
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

ArchiveConfig{
    LogFileDir:    "", // Default will change this "" to Working directory
    ErrorFileName: "error.log",
    WarnFileName:  "warn.log",
    InfoFileName:  "info.log",
    DebugFileName: "debug.log",
    MaxSize:       11,
    MaxBackups:    10,
    MaxAge:        30,
    Compress:      false,
}

```

When you already `New()` that config was build a `logger` register to `Zap`
you can use `zap.L()` get this `logger` by `zap` or get this by this package `log.Get()`

```go
import (
"github.uptutu/log"
)

func main() {
conf := log.NewDefaultConfig()
conf.ArchConf = nil
log.NewLogger(conf)

log.Get().Info("xxxxx")
log.Get().Debug("xxxxx)
}
```