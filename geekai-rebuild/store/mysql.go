package store

import (
	"time"

	"geekai-rebuild/core/types"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewGormConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "geekai_",
			SingularTable: false,
		},
	}
}

func NewMysql(gormConfig *gorm.Config, appConfig *types.AppConfig, log *zap.SugaredLogger) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(appConfig.MysqlDsn), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(32)
	sqlDB.SetMaxOpenConns(512)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Infow("mysql connected", "dsn", appConfig.MysqlDsn)
	return db, nil
}
