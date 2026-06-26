package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/MiKaMoRe/medical-task-tracker/internal/config"
)

func New(env string) (Logger, error) {
	level := slog.LevelInfo
	if env == config.EnvDev {
		level = slog.LevelDebug
	}
	return NewWithLevel(env, level)
}

func NewWithConfigLevel(env, level string) (Logger, error) {
	return NewWithLevel(env, parseSlogLevel(level))
}

func NewWithLevel(env string, level slog.Level) (Logger, error) {
	var (
		l   Logger
		err error
	)

	switch env {
	case config.EnvDev:
		l, err = newDevLogger(level)
	case config.EnvTest:
		l = newJSONLogger(os.Stdout, level)
	case config.EnvProd:
		l = newJSONLogger(os.Stdout, level)
	default:
		return nil, fmt.Errorf("logger: unknown environment: %q", env)
	}

	if err != nil {
		return nil, err
	}

	l.Info("logger initialized", "env", env)
	return l, nil
}

func MustNew(env string) Logger {
	l, err := New(env)
	if err != nil {
		panic(err)
	}
	return l
}

func MustNewWithConfigLevel(env, level string) Logger {
	l, err := NewWithConfigLevel(env, level)
	if err != nil {
		panic(err)
	}
	return l
}
