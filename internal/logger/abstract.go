package logger

import (
	"log/slog"
)

type AppLogger interface {
	Info(message string, args ...slog.Attr)
	Error(message string, err error, args ...slog.Attr)
	Fatal(message string, err error, args ...slog.Attr)
	With(args ...slog.Attr) AppLogger
}
