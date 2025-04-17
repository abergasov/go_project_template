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

		fetchedCh1 := make([]string, 0, len(messages))
		f1Mutex := sync.Mutex{}
		fetchedCh2 := make([]string, 0, len(messages))
		f2Mutex := sync.Mutex{}
		go func() {
			for m := range chBr1 {
				f1Mutex.Lock()
				fetchedCh1 = append(fetchedCh1, m)
				f1Mutex.Unlock()
			}
		}()
		go func() {
			for m := range chBr2 {
				f2Mutex.Lock()
				fetchedCh2 = append(fetchedCh2, m)
				f2Mutex.Unlock()
			}
		}()
		require.Eventually(t, func() bool {
			valid1 := false
			f1Mutex.Lock()
			valid1 = len(fetchedCh1) == len(messages)
			f1Mutex.Unlock()
			valid2 := false
			f2Mutex.Lock()
			valid2 = len(fetchedCh2) == len(messages)
			f2Mutex.Unlock()
			return valid1 && valid2
		}, 10*time.Second, 1*time.Second, "expected all messages to be received")
		require.Equal(t, messages, fetchedCh1, "listener l1: expected %s, got %s", messages, fetchedCh1)
		require.Equal(t, messages, fetchedCh2, "listener l2: expected %s, got %s", messages, fetchedCh2)
	})
	t.Run("check no listeners", func(t *testing.T) {
		br := utils.NewBroadcaster[int]()
		br.Broadcast(100)
	})
}
