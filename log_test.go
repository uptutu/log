package log

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestNew(t *testing.T) {
	assert.NotEqual(t, zap.L(), zap.NewNop())
	assert.NotNil(t, New())
}

func TestNewLogger(t *testing.T) {
	conf := NewDefaultConfig()
	conf.ArchConf = nil
	assert.NotNil(t, NewLogger(conf))
}

func TestGet(t *testing.T) {
	New()
	assert.Equal(t, zap.L(), Get())
}

func TestSet(t *testing.T) {
	log := &zap.Logger{}
	Set(log)
	assert.Equal(t, zap.L(), log)
}

func TestSetLogFields(t *testing.T) {
	fields := map[string]string{
		"source": "test",
	}
	SetLogWithFields(fields)
	Info("try")
}
