package mysql

import (
	"fmt"
	"time"

	"github.com/template-be-commerce/config"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB opens a GORM *gorm.DB backed by MySQL and configures the
// underlying connection pool according to cfg.
func NewDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	logLevel := logger.Warn
	if cfg.MaxOpenConns == 0 { // treat as dev when pool is unconfigured
		logLevel = logger.Info
	}

	db, err := gorm.Open(gmysql.Open(cfg.DSN()), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		return nil, fmt.Errorf("gorm open mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("mysql ping: %w", err)
	}

	return db, nil
}
