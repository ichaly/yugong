package base

import (
	"fmt"
	"github.com/ichaly/yugong/core/data"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnect(che Cache, cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(buildDialect(cfg.Database), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return nil, err
	}
	err = db.Use(che)
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(
		&data.Author{},
		&data.Video{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func buildDialect(ds *DataSource) gorm.Dialector {
	args := []interface{}{ds.Username, ds.Password, ds.Host, ds.Port, ds.Name}
	if ds.Dialect == "mysql" {
		return mysql.Open(fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", args...,
		))
	} else {
		return postgres.Open(fmt.Sprintf(
			"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable TimeZone=Asia/Shanghai", args...,
		))
	}
}
