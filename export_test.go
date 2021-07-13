package log

import (
	"testing"
)

func TestInfo(t *testing.T) {
	New()
	Info("test")

	log := WrapFields(map[string]string{
		"from": "test",
	})

	Set(log)

	Info("test for fields is valid")
}

func TestDebug(t *testing.T) {
	New()
	Debug("test debug")
}

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	New()
	Panic("test panic")
}
