package logger

import (
	"log/slog"
	"strings"
)

type StringWith struct {
	Key string
	Val string
}

func (s StringWith) slog() slog.Attr {
	key := s.Key
	if !strings.HasPrefix(key, "_") {
		key = "_" + key
	}
	return slog.String(key, s.Val)
}

type AppLogger interface {
	Info(message string, args ...StringWith)
	Error(message string, err error, args ...StringWith)
	Fatal(message string, err error, args ...StringWith)
	With(args ...StringWith) AppLogger
}

func WithString(key, val string) StringWith {
	return StringWith{Key: key, Val: val}
}
