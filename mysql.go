package gormxmysql

import (
	"errors"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	gormx "c5x.io/data-gorm"
	"c5x.io/logx"
)

var log = logx.New().Category("data-gorm").Component("mysql")

var (
	connOnce sync.Once
)

func init() {
	gormx.RegisterDriver(gormx.DriverMysql, new(MysqlProvider))
}

type MysqlProvider struct {
}

//Connect impl DatabaseProvider for gorm postgres
func (p MysqlProvider) Connect(config *gormx.DatabaseConfig) (*gorm.DB, error) {
	if config.Dialect == gormx.DriverMysql {
		if db, err := gorm.Open(mysql.New(mysql.Config{DSN: config.DSN}), &gorm.Config{
			Logger: gormx.DefaultLogger(&config.Logger),
		}); err == nil {
			if sqlDB, err := db.DB(); err == nil {
				if config.MaxIdle > 0 {
					sqlDB.SetMaxIdleConns(config.MaxIdle)
				}
				if config.MaxOpen > 0 && config.MaxOpen > config.MaxIdle {
					sqlDB.SetMaxOpenConns(100)
				}
				if config.MaxLifetime > 0 {
					sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)
				}
				return db, nil
			} else {
				return nil, errors.New("open DB failed")
			}
		} else {
			log.Errorf("connect mysql db failed: error=%s", err.Error())
		}
		return nil, errors.New("connect db failed")
	}
	return nil, errors.New("driver is not postgres")
}
