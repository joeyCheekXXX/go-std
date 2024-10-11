package mysql

import (
	"github.com/joeyCheekXXX/go-std/db/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	config.GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (m *Config) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}

// NewMysql 初始化Mysql数据库
func NewMysql(m *Config) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), config.Gorm.Config(&m.GeneralDB)); err != nil {
		return nil
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
