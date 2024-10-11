package test

import (
	"github.com/joeyCheekXXX/go-std/log"
	"testing"
)

func TestLog(t *testing.T) {

	log.LoadDriver(&log.Config{
		Level: "debug",
	})
	log.Logger.Debug("hello world")
	log.Logger.Debug("hello world", log.String("key", "value"))
}
