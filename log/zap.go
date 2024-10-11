package log

import (
	"github.com/joeyCheekXXX/go-std/directory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type _zapCore struct {
	level      zapcore.Level
	encoder    zapcore.Encoder
	writerSync *lumberjack.Logger
	zapcore.Core
}

func newZapCore(config *Config, level zapcore.Level) *_zapCore {
	entity := &_zapCore{
		level:   level,
		encoder: config.Encoder(),
		writerSync: &lumberjack.Logger{
			Filename:   config.WriterSync.Director + "/" + config.Prefix + "/" + level.String() + ".log",
			MaxSize:    config.WriterSync.MaxSize,
			MaxAge:     config.WriterSync.MaxAge,
			MaxBackups: config.WriterSync.MaxBackups,
			LocalTime:  true,
			Compress:   config.WriterSync.Compress,
		},
		Core: nil,
	}

	syncer := zapcore.AddSync(os.Stdout)
	if config.EnableWire {
		// 判断是否有Director文件夹
		if ok, _ := directory.PathExists(config.WriterSync.Director); !ok {
			_ = os.Mkdir(config.WriterSync.Director, os.ModePerm)
		}
		syncer = entity.WriteSyncer()
	}

	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == level
	})
	entity.Core = zapcore.NewCore(config.Encoder(), syncer, levelEnabler)
	return entity
}

func (z *_zapCore) WriteSyncer() zapcore.WriteSyncer {
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(z.writerSync), zapcore.AddSync(os.Stdout))
}

func (z *_zapCore) Enabled(level zapcore.Level) bool {
	return z.level == level
}

func (z *_zapCore) With(fields []zapcore.Field) zapcore.Core {
	return z.Core.With(fields)
}

func (z *_zapCore) Check(entry zapcore.Entry, check *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if z.Enabled(entry.Level) {
		return check.AddCore(entry, z)
	}
	return check
}

func (z *_zapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	for i := 0; i < len(fields); i++ {
		if fields[i].Key == "business" || fields[i].Key == "folder" || fields[i].Key == "directory" {
			z.Core = zapcore.NewCore(z.encoder, z.WriteSyncer(), z.level)
		}
	}
	return z.Core.Write(entry, fields)
}

func (z *_zapCore) Sync() error {
	return z.Core.Sync()
}
