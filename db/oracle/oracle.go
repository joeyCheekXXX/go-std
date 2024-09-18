package oracle

import (
	"github.com/joeyCheek888/go-std/db/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	config.GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (m *Config) Dsn() string {
	return "oracle://" + m.Username + ":" + m.Password + "@" + m.Path + ":" + m.Port + "/" + m.Dbname + "?" + m.Config

}

// NewOracle 初始化oracle数据库
// 如果需要Oracle库 放开import里的注释 把下方 mysql.Config 改为 oracle.Config ;  mysql.New 改为 oracle.New
func NewOracle(m *Config) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}
	oracleConfig := mysql.Config{
		DSN:               m.Dsn(), // DSN data source name
		DefaultStringSize: 191,     // string 类型字段的默认长度
	}
	if db, err := gorm.Open(mysql.New(oracleConfig), config.Gorm.Config(&m.GeneralDB)); err != nil {
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
