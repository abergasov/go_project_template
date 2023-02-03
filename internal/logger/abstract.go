package logger

import "go.uber.org/zap/zapcore"

type AppLogger interface {
	Info(message string, args ...zapcore.Field)
	Error(message string, err error, args ...zapcore.Field)
	Fatal(message string, err error, args ...zapcore.Field)
	With(arg zapcore.Field) AppLogger
}
