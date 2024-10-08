package httpx

import (
	"github.com/joeyCheek888/go-std/log"
	"go.uber.org/zap"
)

type Config struct {
	Addr             string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Port             string `mapstructure:"port" json:"port" yaml:"port"`
	RouterPrefix     string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
	Mode             string `mapstructure:"mode" json:"mode" yaml:"mode"`                      // debug release test
	EnableSwaggerDoc bool   `mapstructure:"swagger-doc" json:"swagger-doc" yaml:"swagger-doc"` // 是否启用swagger文档
}

func (c *Config) Check() {

	if c.Addr == "" {
		c.Addr = "0.0.0.0"
	}

	if c.Port == "" {
		c.Port = "8080"
	}

	if c.Mode == "" {
		c.Mode = "test"
	}

	log.Logger.Info("http配置", zap.Any("config", c))
}
