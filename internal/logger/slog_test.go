package logger_test

import (
	"fmt"
	"go_project_template/internal/logger"
	"sync"
	"testing"

	"github.com/google/uuid"
)

func Test_SLogger_purelog_with_stdout(t *testing.T) {
	// given
	appLog := logger.NewAppSLogger("test")
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
		l := logger.NewAppSLogger("test_2")
		l.Error("source log", fmt.Errorf("123"))
	})
}

func Test_DefaultLogger(t *testing.T) {
	// given
	appLog := logger.NewAppSLogger("test")
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
	appLog := logger.NewAppSLogger("test")

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
