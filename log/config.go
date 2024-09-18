package log

import (
	"go.uber.org/zap/zapcore"
	"time"
)

// Config
// @Description: 主配置
type Config struct {
	Level         string     `mapstructure:"level" json:"level" yaml:"level"`                            // 级别 debug info warn error
	Prefix        string     `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 日志前缀 log
	Format        string     `mapstructure:"format" json:"format" yaml:"format"`                         // 输出格式 json console
	EncodeLevel   string     `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`       // 编码级
	StacktraceKey string     `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"` // 栈名
	ShowLine      bool       `mapstructure:"show-line" json:"show-line" yaml:"show-line"`                // 显示行
	EnableWire    bool       `mapstructure:"enable-wire" json:"enable-wire" yaml:"enable-wire"`          // 是否开启写入到文件夹
	WriterSync    WriterSync `mapstructure:"writer-sync" json:"writer-sync" yaml:"writer-sync"`          // 日志文件配置
	EnablePush    bool       `mapstructure:"enable-push" json:"enable-push" yaml:"enable-push"`          // 是否开启推送
}

func initDefault() *Config {
	return &Config{
		Level:         "info",
		Prefix:        "app",
		Format:        "console",
		EncodeLevel:   "LowercaseColorLevelEncoder",
		StacktraceKey: "stacktrace",
		ShowLine:      true,
		EnableWire:    false,
		WriterSync:    WriterSync{},
		EnablePush:    false,
	}
}

// check
//
//	@Description: 配置检查
//	@receiver c *Config
func (c *Config) check() {
	if c.Level == "" {
		c.Level = "info"
	}

	if c.Prefix == "" {
		c.Prefix = "app"
	}

	if c.Format == "" {
		c.Format = "console"
	}

	if c.EncodeLevel == "" {
		c.EncodeLevel = "LowercaseColorLevelEncoder"
	}

	if c.StacktraceKey == "" {
		c.StacktraceKey = "stacktrace"
	}

	if c.EnableWire {

		if c.WriterSync.Director == "" {
			c.WriterSync.Director = "log"
		}

		if c.WriterSync.MaxSize == 0 {
			c.WriterSync.MaxSize = 10
		}

	}
}

// WriterSync
// @Description: 写入配置
type WriterSync struct {
	Director   string `mapstructure:"director" json:"director" yaml:"director"`          // 日志文件的位置，也就是路径
	MaxSize    int    `mapstructure:"max-size" json:"max-size" yaml:"max-size"`          // 在进行切割之前，日志文件的最大大小（以MB为单位）
	MaxBackups int    `mapstructure:"max-backups" json:"max-backups" yaml:"max-backups"` // 保留旧文件的最大个数
	MaxAge     int    `mapstructure:"max-age" json:"max-age" yaml:"max-age"`             // 保留旧文件的最大天数
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`          // 是否压缩/归档旧文件
}

// Levels 根据字符串转化为 zapcore.Levels
func (c *Config) Levels() []zapcore.Level {
	levels := make([]zapcore.Level, 0, 7)
	level, err := zapcore.ParseLevel(c.Level)
	if err != nil {
		level = zapcore.DebugLevel
	}
	for ; level <= zapcore.FatalLevel; level++ {
		levels = append(levels, level)
	}
	return levels
}

func (c *Config) Encoder() zapcore.Encoder {
	config := zapcore.EncoderConfig{
		TimeKey:       "time",
		NameKey:       "name",
		LevelKey:      "level",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: c.StacktraceKey,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString("[" + c.Prefix + "]" + t.Format("2006-01-02 15:04:05.000"))
		},
		EncodeLevel:    c.LevelEncoder(),
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
	if c.Format == "json" {
		return zapcore.NewJSONEncoder(config)
	}
	return zapcore.NewConsoleEncoder(config)

}

// LevelEncoder 根据 EncodeLevel 返回 zapcore.LevelEncoder
func (c *Config) LevelEncoder() zapcore.LevelEncoder {
	switch {
	case c.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case c.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case c.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case c.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}
