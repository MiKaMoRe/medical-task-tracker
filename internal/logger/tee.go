package logger

import (
	"context"
	"errors"
	"log/slog"
)

type teeHandler struct {
	a, b slog.Handler
}

func newTeeHandler(a, b slog.Handler) *teeHandler {
	return &teeHandler{a, b}
}

func (t *teeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return t.a.Enabled(ctx, level) || t.b.Enabled(ctx, level)
}

func (t *teeHandler) Handle(ctx context.Context, r slog.Record) error {
	err1 := t.a.Handle(ctx, r)
	err2 := t.b.Handle(ctx, r.Clone())
	return errors.Join(err1, err2)
}

func (t *teeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &teeHandler{
		a: t.a.WithAttrs(attrs),
		b: t.b.WithAttrs(attrs),
	}
}

func (t *teeHandler) WithGroup(name string) slog.Handler {
	return &teeHandler{
		a: t.a.WithGroup(name),
		b: t.b.WithGroup(name),
	}
}
