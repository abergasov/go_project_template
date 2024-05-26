package logger_test

import (
	"bytes"
	"fmt"
	"go_project_template/internal/logger"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
)

func Test_SLogger_purelog_with_stdout(t *testing.T) {
	// given
	appLog := newTestLogger(t)
	appLog.Info("test")
	appLog.Error("source log", fmt.Errorf("123"))

	// when, then
	concurrentlyLogIt(
		appLog.With(
			logger.WithString("a", "b"),
			logger.WithString("c", "d"),
		),
	)
	t.Run("pure check", func(t *testing.T) {
		l := newTestLogger(t)
		l.Error("source log", fmt.Errorf("123"))
	})
}

func Test_DefaultLogger(t *testing.T) {
	// given
	appLog := newTestLogger(t)
	appLog.Info("test")
	appLog.Error("source log", fmt.Errorf("123"))

	// when, then
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			concurrentlyLogIt(
				appLog.With(
					logger.WithString("a", "b"),
					logger.WithString("c", "d"),
				),
			)
		}()
		go func() {
			defer wg.Done()
			defaultAppLog := logger.NewAppSLogger("test_2")
			concurrentlyLogIt(
				defaultAppLog.With(
					logger.WithString("a", "b"),
					logger.WithString("c", "d"),
				),
			)
		}()
	}
	wg.Wait()
}

func Test_SLogger_multiply_writers(t *testing.T) {
	// given
	appLog := newTestLogger(t)

	// when, then
	concurrentlyLogIt(appLog)

	t.Run("pure check", func(t *testing.T) {
		l := logger.NewAppSLogger(
			"test_2",
			logger.WithString("additional", "value"),
			logger.WithString("additional2", "value"),
		)
		l.Info("test")
	})
}

func concurrentlyLogIt(appLog logger.AppLogger) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				appLog.Info("message")
				appLog.Error("message", fmt.Errorf("error"))
				appLog.Error("message", fmt.Errorf("error"), logger.WithString("key", uuid.NewString()))
			}
		}()
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			customLogger1 := appLog.With(logger.WithString("keyA", uuid.NewString()))
			customLogger1.Info("customLogger message")
			customLogger1.Error("customLogger message", fmt.Errorf("customLogger error"))
			customLogger1.Error("customLogger message", fmt.Errorf("customLogger error"), logger.WithString("keyB", uuid.NewString()))
		}()
	}
	wg.Wait()
}

type TestLogger struct {
	logsS      [][]byte
	logs       bytes.Buffer
	std        io.Writer
	logsChan   chan []byte
	signalChan chan struct{}
	wg         sync.WaitGroup
}

func (tl *TestLogger) Write(p []byte) (n int, err error) {
	b := make([]byte, len(p)) // slog reuse the buffer, so we need to copy it to avoid race condition
	copy(b, p)
	tl.logsChan <- b
	return 0, nil
}

func (tl *TestLogger) process() {
	defer tl.wg.Done()
	for {
		select {
		case <-tl.signalChan:
			return
		case data := <-tl.logsChan:
			tl.logs.Write(data)
		}
	}
}

func newTestLogger(t testing.TB) logger.AppLogger {
	tl := &TestLogger{
		std:        os.Stdout,
		logsChan:   make(chan []byte, 1_000),
		signalChan: make(chan struct{}),
	}
	tl.wg.Add(1)
	go tl.process()
	t.Cleanup(func() {
		if t.Failed() {
			tl.signalChan <- struct{}{}
			tl.wg.Wait()
			print(tl.logs.String())
		}
	})
	return logger.InitLogger([]io.Writer{tl}, "test")
}
