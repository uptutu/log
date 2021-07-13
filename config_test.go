package log

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"testing"
)

func TestDefaultArchiveConfig(t *testing.T) {
	assert.NotNil(t, DefaultArchiveConfig())
	assert.Nil(t, DefaultArchiveConfig().errLog)
	assert.NotNil(t, DefaultArchiveConfig().ErrorFileName)
}

func TestNewDefaultWithoutArchiveConfig(t *testing.T) {
	assert.Nil(t, NewDefaultWithoutArchiveConfig().ArchConf)
}

func TestDefaultZapEncoderConfig(t *testing.T) {
	assert.NotNil(t, DefaultZapEncoderConfig())
	assert.Equal(t, "content", DefaultZapEncoderConfig().MessageKey)
}

func TestNewDefaultConfig(t *testing.T) {
	assert.True(t, NewDefaultConfig().Stdout)
	assert.Equal(t, *DefaultArchiveConfig(), *NewDefaultConfig().ArchConf)
}

func TestConstVar(t *testing.T) {
	assert.Equal(t, 1, int(WithArchive))
	assert.Equal(t, 2, int(WithOutArchive))
	assert.Equal(t, 255, int(ErrConfig))
}

func TestConfigCheckWithArchConfLogFileDirIsEmptyString(t *testing.T) {
	ac := DefaultArchiveConfig()
	conf := Config{ArchConf: ac}
	assert.Equal(t, WithArchive, conf.check())
	wd, err := filepath.Abs(filepath.Dir(filepath.Join(".")))
	assert.Nil(t, err)
	assert.Equal(t, wd+"/logs", conf.ArchConf.LogFileDir)
}

func TestConfigCheckWithArchConfLogFileDir(t *testing.T) {
	ac := DefaultArchiveConfig()
	ac.LogFileDir = "/test/logs"
	conf := Config{ArchConf: ac}
	assert.Equal(t, WithArchive, conf.check())
	assert.Equal(t, "/test/logs", conf.ArchConf.LogFileDir)

	conf.ArchConf.LogFileDir = "test/logs/"
	assert.Equal(t, WithArchive, conf.check())
	assert.Equal(t, "test/logs", conf.ArchConf.LogFileDir)
}

func TestConfigCheckWithErrConfig(t *testing.T) {
	ac := ArchiveConfig{
		LogFileDir:    "",
		ErrorFileName: "",
		WarnFileName:  "",
		InfoFileName:  "",
		DebugFileName: "",
		MaxSize:       30,
	}
	conf := Config{ArchConf: &ac}
	assert.Equal(t, ErrConfig, conf.check())
}

func TestConfigCheckWithErrFileNameConfig(t *testing.T) {
	ac := ArchiveConfig{
		LogFileDir:    "",
		ErrorFileName: "/logs.log/",
		MaxSize:       30,
	}
	conf := Config{ArchConf: &ac}
	assert.Equal(t, WithArchive, conf.check())
	assert.Equal(t, "logs.log", conf.ArchConf.ErrorFileName)
}

func TestConfigBuildArchLogger(t *testing.T) {
	conf := NewDefaultConfig()
	assert.Nil(t, conf.ArchConf.infoLog)
	assert.Nil(t, conf.ArchConf.errLog)
	assert.Nil(t, conf.ArchConf.debugLog)
	assert.Nil(t, conf.ArchConf.warnLog)
	conf.buildArchLogger()
	assert.NotNil(t, conf.ArchConf.infoLog)
	assert.NotNil(t, conf.ArchConf.errLog)
	assert.NotNil(t, conf.ArchConf.debugLog)
	assert.NotNil(t, conf.ArchConf.warnLog)
}

func TestWrapSyncOpts(t *testing.T) {
	conf := NewDefaultConfig()
	conf.Stdout = false
	conf.buildArchLogger()
	opts := conf.wrapSyncOpts()
	assert.NotNil(t, opts)
}

func TestConfigBuild(t *testing.T) {
	conf := NewDefaultConfig()
	conf.ArchConf = &ArchiveConfig{MaxSize: 30}
	log, err := conf.build()
	assert.NotNil(t, err)
	assert.Nil(t, log)
	assert.Equal(t, errors.New("log config error"), err)

	conf = NewDefaultConfig()
	conf.ArchConf = nil
	log, err = conf.build()
	assert.Nil(t, err)
	assert.NotNil(t, log)

	conf = NewDefaultConfig()
	log, err = conf.build()
	assert.Nil(t, err)
	assert.NotNil(t, log)
}

func TestConfig_SetLevel(t *testing.T) {
	conf := NewDefaultConfig()
	assert.Equal(t, zapcore.DebugLevel, conf.Level.Level())
	conf.SetLevel("info")
	assert.Equal(t, zapcore.InfoLevel, conf.Level.Level())
	conf.SetLevel("debug")
	assert.Equal(t, zapcore.DebugLevel, conf.Level.Level())
	conf.SetLevel("warn")
	assert.Equal(t, zapcore.WarnLevel, conf.Level.Level())
	conf.SetLevel("error")
	assert.Equal(t, zapcore.ErrorLevel, conf.Level.Level())
	conf.SetLevel("panic")
	assert.Equal(t, zapcore.PanicLevel, conf.Level.Level())
	conf.SetLevel("fatal")
	assert.Equal(t, zapcore.FatalLevel, conf.Level.Level())
	conf.SetLevel("xxx")
	assert.Equal(t, zapcore.InfoLevel, conf.Level.Level())
}
