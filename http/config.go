package http

type Config struct {
	Addr             string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Port             string `mapstructure:"port" json:"port" yaml:"port"`
	RouterPrefix     string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
	Mode             string `mapstructure:"mode" json:"mode" yaml:"mode"` // debug release
	LimitCountIP     int    `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"`
	LimitTimeIP      int    `mapstructure:"iplimit-time" json:"iplimit-time" yaml:"iplimit-time"`
	UseMultipoint    bool   `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"` // 多点登录拦截
	EnableSwaggerDoc bool   `mapstructure:"swagger-doc" json:"swagger-doc" yaml:"swagger-doc"`
}
