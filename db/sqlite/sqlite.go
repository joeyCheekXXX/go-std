package sqlite

import (
	"github.com/glebarez/sqlite"
	"go-std/db/config"
	"gorm.io/gorm"
	"path/filepath"
)

type Config struct {
	config.GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (s *Config) Dsn() string {
	return filepath.Join(s.Path, s.Dbname+".db")
}

// NewSqlite 初始化Sqlite数据库
func NewSqlite(m *Config) *gorm.DB {

	if m.Dbname == "" {
		return nil
	}

	if db, err := gorm.Open(sqlite.Open(m.Dsn()), config.Gorm.Config(&m.GeneralDB)); err != nil {
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
