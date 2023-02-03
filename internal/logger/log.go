package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	l *zap.Logger
}

func NewAppLogger(appHash string) (*Logger, error) {
	cnf := zap.NewProductionConfig()
	cnf.DisableStacktrace = true
	cnf.DisableCaller = true
	cnf.EncoderConfig.TimeKey = zapcore.OmitKey

	z, err := cnf.Build()
	if err != nil {
		return nil, err
	}

	if appHash != "" {
		z = z.With(zap.String("hash", appHash))
	}
	return &Logger{l: z}, nil
}

func (a Logger) Info(message string, args ...zapcore.Field) {
	a.l.Info(message, args...)
}

func (a Logger) Error(message string, err error, args ...zapcore.Field) {
	if len(args) == 0 {
		a.l.Error(message, zap.Error(err))
		return
	}
	a.l.Error(message, prepareParams(err, args)...)
}

func (a Logger) Fatal(message string, err error, args ...zapcore.Field) {
	if len(args) == 0 {
		a.l.Fatal(message, zap.Error(err))
		return
	}
	a.l.Fatal(message, prepareParams(err, args)...)
}

func (a Logger) With(arg zapcore.Field) AppLogger {
	return Logger{l: a.l.With(arg)}
}

func prepareParams(err error, args []zapcore.Field) []zapcore.Field {
	params := make([]zapcore.Field, 0, len(args)+1)
	params = append(params, zap.Error(err))
	params = append(params, args...)
	return params
}
