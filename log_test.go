package log

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestNew(t *testing.T) {
	assert.Equal(t, New(), zap.L())
}

func TestNewLogger(t *testing.T) {
	conf := NewDefaultConfig()
	conf.ArchConf = nil
	assert.Equal(t, NewLogger(conf), zap.L())
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
