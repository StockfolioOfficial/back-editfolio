package di

import (
	"github.com/stockfolioofficial/back-editfolio/core/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/go-sql-driver/mysql"
)

func NewDatabase() (db *gorm.DB) {
	var logLevel = logger.Info

	if !config.IsDebug {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(mysql.Open(config.DBConn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	return
}