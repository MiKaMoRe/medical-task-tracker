package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/MiKaMoRe/medical-task-tracker/internal/config"
)

type slogLogger struct {
	inner   *slog.Logger
	cleanup func() error
}

func (l *slogLogger) Debug(msg string, args ...any) { l.inner.Debug(msg, args...) }
func (l *slogLogger) Info(msg string, args ...any)  { l.inner.Info(msg, args...) }
func (l *slogLogger) Warn(msg string, args ...any)  { l.inner.Warn(msg, args...) }
func (l *slogLogger) Error(msg string, args ...any) { l.inner.Error(msg, args...) }

func (l *slogLogger) Close() error {
	if l.cleanup != nil {
		return l.cleanup()
	}
	return nil
}

func newDevLogger(level slog.Level) (Logger, error) {
	stdout := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})

	f, err := openLogFile()
	if err != nil || f == nil {
		return &slogLogger{inner: slog.New(stdout)}, nil
	}

	file := slog.NewJSONHandler(f, &slog.HandlerOptions{Level: level})
	h := newTeeHandler(stdout, file)
	return &slogLogger{
		inner:   slog.New(h),
		cleanup: f.Close,
	}, nil
}

func newJSONLogger(w io.Writer, level slog.Level) Logger {
	h := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})
	return &slogLogger{inner: slog.New(h)}
}

func parseSlogLevel(level string) slog.Level {
	switch level {
	case config.LogLevelDebug:
		return slog.LevelDebug
	case config.LogLevelWarn:
		return slog.LevelWarn
	case config.LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
