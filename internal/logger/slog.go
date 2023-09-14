package logger

import (
	"log/slog"
	"os"
)

type SLogger struct {
	l *slog.Logger
}

var _ AppLogger = (*SLogger)(nil)

func NewAppSLogger(appHash string) *SLogger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if appHash != "" {
		logger = logger.With(slog.String("hash", appHash))
	}

	return &SLogger{l: logger}
}

func (l *SLogger) Info(message string, args ...slog.Attr) {
	l.l.Info(message, prepareSlogParams(nil, args)...)
}

func (l *SLogger) Error(message string, err error, args ...slog.Attr) {
	l.l.Error(message, prepareSlogParams(err, args)...)
}

func (l *SLogger) Fatal(message string, err error, args ...slog.Attr) {
	l.l.Error(message, prepareSlogParams(err, args)...)
	os.Exit(1)
}

func (l *SLogger) With(args ...slog.Attr) AppLogger {
	return &SLogger{
		l: l.l.With(prepareSlogParams(nil, args)...),
	}
}

func prepareSlogParams(err error, args []slog.Attr) []any {
	params := make([]any, 0, len(args)+1)
	if err != nil {
		params = append(params, err)
	}
	for _, arg := range args {
		params = append(params, arg)
	}
	return params
}
