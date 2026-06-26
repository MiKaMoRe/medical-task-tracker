package gormdb

import (
	"context"
	"fmt"
	"time"

	"github.com/MiKaMoRe/medical-task-tracker/internal/httpctx"
	applogger "github.com/MiKaMoRe/medical-task-tracker/internal/logger"
	gormlogger "gorm.io/gorm/logger"
)

const defaultSlowQueryThreshold = 500 * time.Millisecond

type structuredGormLogger struct {
	appLogger     applogger.Logger
	level         gormlogger.LogLevel
	slowThreshold time.Duration
}

func newStructuredGormLogger(appLogger applogger.Logger, level gormlogger.LogLevel) gormlogger.Interface {
	if level == 0 {
		level = gormlogger.Error
	}

	return &structuredGormLogger{
		appLogger:     appLogger,
		level:         level,
		slowThreshold: defaultSlowQueryThreshold,
	}
}

func (l *structuredGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	cloned := *l
	cloned.level = level
	return &cloned
}

func (l *structuredGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level < gormlogger.Info {
		return
	}
	l.appLogger.Info(
		"db info",
		"message", fmt.Sprintf(msg, data...),
		"requestId", requestIDFromContext(ctx),
	)
}

func (l *structuredGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level < gormlogger.Warn {
		return
	}
	l.appLogger.Warn(
		"db warn",
		"message", fmt.Sprintf(msg, data...),
		"requestId", requestIDFromContext(ctx),
	)
}

func (l *structuredGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level < gormlogger.Error {
		return
	}
	l.appLogger.Error(
		"db error",
		"message", fmt.Sprintf(msg, data...),
		"requestId", requestIDFromContext(ctx),
	)
}

func (l *structuredGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level == gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)

	shouldLogInfo := l.level >= gormlogger.Info
	shouldLogSlow := l.level >= gormlogger.Warn && l.slowThreshold > 0 && elapsed > l.slowThreshold
	shouldLogError := l.level >= gormlogger.Error && err != nil
	if !shouldLogInfo && !shouldLogSlow && !shouldLogError {
		return
	}

	sql, rows := fc()
	args := []any{
		"elapsedMs", elapsed.Milliseconds(),
		"rows", rows,
		"sql", sql,
		"requestId", requestIDFromContext(ctx),
	}

	switch {
	case shouldLogError:
		l.appLogger.Error("db query failed", append(args, "error", err.Error())...)
	case shouldLogSlow:
		l.appLogger.Warn("db slow query", args...)
	case shouldLogInfo:
		l.appLogger.Info("db query executed", args...)
	}
}

func requestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return "unknown"
	}
	if id, ok := httpctx.RequestID(ctx); ok {
		return id
	}
	return "unknown"
}
