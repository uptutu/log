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
Edit the config as you want, `Config` field is `uber-go/Zap` config, ArchConf as you want, if set `nil` to this field then
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