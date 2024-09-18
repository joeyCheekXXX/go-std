package mssql

import (
	"github.com/joeyCheek888/go-std.git/db/config"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Config struct {
	config.GeneralDB `yaml:",inline" mapstructure:",squash"`
}

// Dsn "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
func (m *Config) Dsn() string {
	return "sqlserver://" + m.Username + ":" + m.Password + "@" + m.Path + ":" + m.Port + "?database=" + m.Dbname + "&encrypt=disable"
}

// NewMssql 初始化Mssql数据库
func NewMssql(m *Config) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}
	mssqlConfig := sqlserver.Config{
		DSN:               m.Dsn(), // DSN data source name
		DefaultStringSize: 191,     // string 类型字段的默认长度
	}
	if db, err := gorm.Open(sqlserver.New(mssqlConfig), config.Gorm.Config(&m.GeneralDB)); err != nil {
		return nil
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
