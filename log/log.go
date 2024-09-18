package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {

	LoadDriver(initDefault())

	Logger.Debug("logger init success ")
}

// LoadDriver
//
//	@Description:
//	@param config
func LoadDriver(config *Config) {

	config.check()

	levels := config.Levels()
	length := len(levels)
	cores := make([]zapcore.Core, 0, length)
	for i := 0; i < length; i++ {
		core := newZapCore(config, levels[i])
		cores = append(cores, core)
	}
	Logger = zap.New(zapcore.NewTee(cores...))
	if config.ShowLine {
		Logger = Logger.WithOptions(zap.AddCaller())
	}

	zap.ReplaceGlobals(Logger)
}
