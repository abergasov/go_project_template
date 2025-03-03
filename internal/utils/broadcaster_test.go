package utils_test

import (
	"go_project_template/internal/utils"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBroadcasterBroadcast(t *testing.T) {
	b := utils.NewBroadcaster[int]()
	ch1 := b.RegisterListener("l1")
	ch2 := b.RegisterListener("l2")

	msg := 42

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		select {
		case m := <-ch1:
			require.Equalf(t, msg, m, "listener l1: expected %d, got %d", msg, m)
		case <-time.After(100 * time.Millisecond):
			t.Error("listener l1 did not receive message")
		}
	}()
	go func() {
		defer wg.Done()
		select {
		case m := <-ch2:
			require.Equalf(t, msg, m, "listener l2: expected %d, got %d", msg, m)
		case <-time.After(100 * time.Millisecond):
			t.Error("listener l2 did not receive message")
		}
	}()
	go b.Broadcast(msg)

	wg.Wait()
	t.Run("check unregister", func(t *testing.T) {
		br := utils.NewBroadcaster[int]()
		ch := br.RegisterListener("l")
		br.UnregisterListener("l")
		_, ok := <-ch
		require.False(t, ok, "broadcaster unregistered listener")
	})
	t.Run("check register", func(t *testing.T) {
		br := utils.NewBroadcaster[string]()
		chBr1 := br.RegisterListener("l1")
		chBr2 := br.RegisterListener("l2")
		messages := []string{"hello", "world", "test"}

		go func() {
			for _, m := range messages {
				br.Broadcast(m)
				time.Sleep(10 * time.Millisecond)
			}
		}()

		for _, expected := range messages {
			select {
			case m := <-chBr1:
				require.Equalf(t, expected, m, "listener l1: expected %s, got %s", expected, m)
			case <-time.After(100 * time.Millisecond):
				t.Errorf("listener l1 did not receive message %s", expected)
			}
			select {
			case m := <-chBr2:
				require.Equalf(t, expected, m, "listener l2: expected %s, got %s", expected, m)
			case <-time.After(100 * time.Millisecond):
				t.Errorf("listener l2 did not receive message %s", expected)
			}
		}
	})
	t.Run("check no listeners", func(t *testing.T) {
		br := utils.NewBroadcaster[int]()
		br.Broadcast(100)
	})
}
