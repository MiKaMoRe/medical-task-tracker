package gormdb

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	LogLevelSilent = logger.Silent
	LogLevelError  = logger.Error
	LogLevelWarn   = logger.Warn
	LogLevelInfo   = logger.Info
)

type Options struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	LogLevel        logger.LogLevel
}

func NewPostgresDialector(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}

func setOptions(opts *Options) {
	if opts.MaxIdleConns == 0 {
		opts.MaxIdleConns = 2
	}

	if opts.MaxOpenConns == 0 {
		opts.MaxOpenConns = 100
	}

	if opts.ConnMaxLifetime == 0 {
		opts.ConnMaxLifetime = 10 * time.Minute
	}

	if opts.ConnMaxIdleTime == 0 {
		opts.ConnMaxIdleTime = 5 * time.Minute
	}

	if opts.LogLevel == 0 {
		opts.LogLevel = logger.Error
	}
}

func Connect(dialector gorm.Dialector, opts *Options) (*gorm.DB, error) {
	setOptions(opts)

	cfgGorm := &gorm.Config{
		Logger: logger.Default.LogMode(opts.LogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(dialector, cfgGorm)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(opts.MaxIdleConns)
	sqlDB.SetMaxOpenConns(opts.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(opts.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(opts.ConnMaxIdleTime)

	return db, nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
