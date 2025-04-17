package logger

import (
	"log/slog"
	"strconv"
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

func WithUnt64(key string, val uint64) StringWith {
	return StringWith{Key: key, Val: strconv.FormatUint(val, 10)}
}

func WithInt64(key string, val int64) StringWith {
	return StringWith{Key: key, Val: strconv.FormatInt(val, 10)}
}

func WithFloat64(key string, val float64) StringWith {
	return StringWith{Key: key, Val: strconv.FormatFloat(val, 'f', -1, 64)}
}

func WithInt(key string, val int) StringWith {
	return StringWith{Key: key, Val: strconv.FormatInt(int64(val), 10)}
}
