package logger

import (
	"context"
	"log/slog"
)

type key struct{}

func ContextWithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, key{}, l)
}

func Get(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(key{}).(*slog.Logger); ok {
		return l
	}

	return slog.Default()
}

func ErrAttr(err error) slog.Attr {
	return slog.Any("err", err)
}
