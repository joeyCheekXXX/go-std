package log

import (
	"testing"
)

func TestLog(t *testing.T) {

	Logger.Info("hello world")
	Logger.Debug("hello world", String("key", "value"))
}
