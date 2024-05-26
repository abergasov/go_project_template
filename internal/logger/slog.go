package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

type SLogger struct {
	logWriters []*slog.Logger // since we not modify this slice, we able avoid mutex usage
}

var _ AppLogger = (*SLogger)(nil)

func NewAppSLogger(appHash string, args ...StringWith) AppLogger {
	return InitLogger([]io.Writer{
		os.Stdout,
	}, appHash, args...)
}

func InitLogger(writers []io.Writer, appHash string, args ...StringWith) AppLogger {
	logs := make([]*slog.Logger, 0, len(writers))
	commit := "unknown"
	if appHash != "" && len(appHash) > 6 {
		commit = appHash[:6]
	}
	for _, w := range writers {
		handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				switch a.Key {
				case "time":
					return slog.Int64("timestamp", time.Now().Unix())
				case "level":
					return slog.String("_level", strings.ToLower(a.Value.String()))
				case "gray_log_level":
					return slog.Int64("level", a.Value.Int64())
				case "msg":
					return slog.String("short_message", a.Value.String())
				}
				return a
			},
		})
		lw := slog.New(handler).With(
			slog.String("version", commit),
		)
		logs = append(logs, lw)
	}
	return &SLogger{logWriters: logs}
}

func (l *SLogger) Info(message string, args ...StringWith) {
	params := prepareSlogParams(nil, args)
	l.processWriters(func(lg *slog.Logger) {
		lg.Info(message, params...)
	})
}

func (l *SLogger) Error(message string, err error, args ...StringWith) {
	params := prepareSlogParams(err, args)
	l.processWriters(func(lg *slog.Logger) {
		lg.Error(message, params...)
	})
}

func (l *SLogger) Fatal(message string, err error, args ...StringWith) {
	params := prepareSlogParams(err, args)
	l.processWriters(func(lg *slog.Logger) {
		lg.Error(message, params...)
	})
	os.Exit(1)
}

func (l *SLogger) With(args ...StringWith) AppLogger {
	logs := make([]*slog.Logger, 0, len(l.logWriters))
	for _, lg := range l.logWriters {
		logs = append(logs, lg.With(prepareSlogParams(nil, args)...))
	}
	return &SLogger{
		logWriters: logs,
	}
}

func prepareSlogParams(err error, args []StringWith) []any {
	params := make([]any, 0, len(args)+2)
	if err != nil {
		params = append(params, WithString("error", err.Error()).slog())
	}
	for _, arg := range args {
		params = append(params, arg.slog())
	}
	return params
}
func (l *SLogger) processWriters(processor func(*slog.Logger)) {
	var wg sync.WaitGroup
	wg.Add(len(l.logWriters))
	for i := range l.logWriters {
		go func(j int) {
			processor(l.logWriters[j])
			wg.Done()
		}(i)
	}
	wg.Wait()
}
