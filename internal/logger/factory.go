package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/MiKaMoRe/medical-task-tracker/internal/config"
)

func New(env string) (Logger, error) {
	var (
		l   Logger
		err error
	)

	switch env {
	case config.EnvDev:
		l, err = newDevLogger()
	case config.EnvTest:
		l = newJSONLogger(os.Stdout, slog.LevelInfo)
	case config.EnvProd:
		l = newJSONLogger(os.Stdout, slog.LevelInfo)
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
