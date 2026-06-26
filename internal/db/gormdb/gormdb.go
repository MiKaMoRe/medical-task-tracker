package gormdb

import (
	"fmt"
	"time"

	applogger "github.com/MiKaMoRe/medical-task-tracker/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const (
	LogLevelSilent = gormlogger.Silent
	LogLevelError  = gormlogger.Error
	LogLevelWarn   = gormlogger.Warn
	LogLevelInfo   = gormlogger.Info
)

type LogLevel = gormlogger.LogLevel

type Options struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	LogLevel        gormlogger.LogLevel
	AppLogger       applogger.Logger
}

func NewPostgresDialector(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}

func setOptions(opts *Options) *Options {
	if opts == nil {
		opts = &Options{}
	}

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
		opts.LogLevel = gormlogger.Error
	}

	return opts
}

func Connect(dialector gorm.Dialector, opts *Options) (*gorm.DB, error) {
	opts = setOptions(opts)

	dbLogger := gormlogger.Default.LogMode(opts.LogLevel)
	if opts.AppLogger != nil {
		dbLogger = newStructuredGormLogger(opts.AppLogger, opts.LogLevel)
	}

	cfgGorm := &gorm.Config{
		Logger: dbLogger,
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
