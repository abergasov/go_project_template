package logger

import (
	"go_project_template/internal/utils"
	"io"
	"log/slog"
	"os"
	"strings"
)

type SLogger struct {
	logger *slog.Logger
}

var _ AppLogger = (*SLogger)(nil)

func NewAppSLogger(args ...Field) AppLogger {
	return InitLogger([]io.Writer{os.Stdout}, args...)
}

func InitLogger(writers []io.Writer, args ...Field) AppLogger {
	handlers := make([]slog.Handler, 0, len(writers))
	for _, w := range writers {
		handlers = append(handlers, slog.NewJSONHandler(w, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				switch a.Key {
				case "level":
					return slog.String("_lvl", strings.ToLower(a.Value.String()))
				}
				return a
			},
		}))
	}
	attrs := make([]any, 0, len(args)+1)
	for _, arg := range args {
		attrs = append(attrs, arg.a)
	}
	attrs = append(attrs, slog.String("commit", utils.GetCommitHash()))

	mh := slog.NewMultiHandler(handlers...)
	return &SLogger{logger: slog.New(mh).With(attrs...)}
}

func (l *SLogger) Info(message string, args ...Field) {
	params := prepareSlogParams(nil, args)
	l.logger.Info(message, params...)
}

func (l *SLogger) Error(message string, err error, args ...Field) {
	params := prepareSlogParams(err, args)
	l.logger.Error(message, params...)
}

func (l *SLogger) Fatal(message string, err error, args ...Field) {
	params := prepareSlogParams(err, args)
	l.logger.Error(message, params...)
	l.logger.Info("fatal error; exiting")
	os.Exit(1)
}

// With creates a child logger with additional structured context.
func (l *SLogger) With(fields ...Field) AppLogger {
	return &SLogger{logger: l.logger.With(prepareSlogParams(nil, fields)...)}
}

func prepareSlogParams(err error, fields []Field) []any {
	params := make([]any, 0, len(fields)+1)
	if err != nil {
		params = append(params, slog.String("error", err.Error()))
	}

	for i := range fields {
		params = append(params, fields[i].a)
	}
	return params
}
