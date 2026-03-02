package logger

import (
	"log/slog"
)

type Field struct {
	a slog.Attr
}

type AppLogger interface {
	Info(message string, args ...Field)
	Error(message string, err error, args ...Field)
	Fatal(message string, err error, args ...Field)
	With(args ...Field) AppLogger
}

func WithString(key, val string) Field {
	return Field{a: slog.String(key, val)}
}

func WithUnt64(key string, val uint64) Field {
	return Field{a: slog.Uint64(key, val)}
}

func WithInt64(key string, val int64) Field {
	return Field{a: slog.Int64(key, val)}
}

func WithFloat64(key string, val float64) Field {
	return Field{a: slog.Float64(key, val)}
}

func WithInt(key string, val int) Field {
	return Field{a: slog.Int(key, val)}
}
