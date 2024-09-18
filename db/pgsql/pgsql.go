package pgsql

import (
	"github.com/joeyCheek888/go-std/db/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	config.GeneralDB `yaml:",inline" mapstructure:",squash"`
}

// Dsn 基于配置文件获取 dsn
func (p *Config) Dsn() string {
	return "host=" + p.Path + " user=" + p.Username + " password=" + p.Password + " dbname=" + p.Dbname + " port=" + p.Port + " " + p.Config
}

// LinkDsn 根据 dbname 生成 dsn
func (p *Config) LinkDsn(dbname string) string {
	return "host=" + p.Path + " user=" + p.Username + " password=" + p.Password + " dbname=" + dbname + " port=" + p.Port + " " + p.Config
}

// NewPgSql 初始化 Postgresql 数据库
func NewPgSql(m *Config) *gorm.DB {

	if m.Dbname == "" {
		return nil
	}
	pgsqlConfig := postgres.Config{
		DriverName:           "",
		DSN:                  m.Dsn(), // DSN data source name
		WithoutQuotingCheck:  false,
		PreferSimpleProtocol: false,
		WithoutReturning:     false,
		Conn:                 nil,
	}
	db, err := gorm.Open(postgres.New(pgsqlConfig), config.Gorm.Config(&m.GeneralDB))
	if err != nil {
		return nil
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	return db
}
