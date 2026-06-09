package utils_test

import (
	"context"
	"go_project_template/internal/utils"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTTLMap(t *testing.T) {
	t.Run("basic put get and expiry", func(t *testing.T) {
		maxTTL := time.Second
		cleanup := time.Second
		m := utils.NewTTLMap[string, int](maxTTL, cleanup)

		m.Put("key", 42)
		require.Equal(t, 1, m.Len())

		val, ok := m.Get("key")
		require.True(t, ok)
		require.Equal(t, 42, val)

		require.Eventually(t, func() bool {
			_, ok = m.Get("key")
			return !ok
		}, 5*cleanup, cleanup)
	})

	t.Run("concurrent access", func(t *testing.T) {
		maxTTL := 2 * time.Second
		cleanup := time.Second
		m := utils.NewTTLMap[string, int](maxTTL, cleanup)
		var wg sync.WaitGroup

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				m.Put("key"+strconv.Itoa(i), i)
			}(i)
		}

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				if v, ok := m.Get("key" + strconv.Itoa(i)); ok {
					require.Equal(t, i, v)
				}
			}(i)
		}
		wg.Wait()
	})

	t.Run("DoAndApply method should modify value", func(t *testing.T) {
		// given
		m := utils.NewTTLMap[string, int](time.Minute, time.Minute)
		m.Put("a", 1)
		val, ok := m.Get("a")
		require.True(t, ok)
		require.Equal(t, 1, val)
		require.False(t, m.DoAndApply("missing", func(v int) int { return v + 1 }))

		// when
		require.True(t, m.DoAndApply("a", func(v int) int {
			v += 10
			return v
		}))

		// then
		val, ok = m.Get("a")
		require.True(t, ok)
		require.Equal(t, 11, val)
	})

	t.Run("do method", func(t *testing.T) {
		m := utils.NewTTLMap[string, int](time.Minute, time.Minute)
		m.Put("a", 1)
		executed := false
		ok := m.Do("a", func(v int) { executed = true; require.Equal(t, 1, v) })
		require.True(t, ok)
		require.True(t, executed)

		ok = m.Do("missing", func(v int) {})
		require.False(t, ok)
	})
	t.Run("upsert method", func(t *testing.T) {
		m := utils.NewTTLMap[string, int](time.Minute, time.Minute)
		m.Upsert("a", func(v int) int { return v + 1 }, 0)
		val, ok := m.Get("a")
		require.True(t, ok)
		require.Equal(t, 0, val)

		m.Upsert("a", func(v int) int { return v + 1 }, 0)
		val, ok = m.Get("a")
		require.True(t, ok)
		require.Equal(t, 1, val)
	})
}

func TestUpsert_RefreshesTTLForExistingKey(t *testing.T) {
	maxTTL := 80 * time.Millisecond
	cleanup := 10 * time.Millisecond
	m := utils.NewTTLMap[string, int](maxTTL, cleanup)
	defer m.Close()

	// create key
	m.Upsert("k", func(v int) int { return v + 1 }, 0)
	v, ok := m.Get("k")
	require.True(t, ok)
	require.Equal(t, 0, v)

	// let it age
	time.Sleep(maxTTL / 2)

	// Upsert should both update value and refresh TTL
	m.Upsert("k", func(v int) int { return v + 1 }, 0)

	// after original TTL but before refreshed TTL, it must still be alive
	time.Sleep(maxTTL / 2)

	v, ok = m.Get("k")
	require.True(t, ok)
	require.Equal(t, 1, v)
	require.Equal(t, 1, len(m.LoadAll()))
}

func TestGetAndRefresh_ExtendsTTL(t *testing.T) {
	maxTTL := 300 * time.Millisecond
	cleanup := 50 * time.Millisecond
	m := utils.NewTTLMap[string, int](maxTTL, cleanup)
	defer m.Close()

	m.Put("x", 7)

	// halfway -> still alive
	time.Sleep(maxTTL / 2)
	v, ok := m.GetAndRefresh("x")
	require.True(t, ok)
	require.Equal(t, 7, v)

	// after initial TTL passed, should still be alive due to refresh
	time.Sleep(maxTTL / 2)
	v, ok = m.Get("x")
	require.True(t, ok)
	require.Equal(t, 7, v)

	// should expire after the refreshed TTL
	require.Eventually(t, func() bool {
		_, ok = m.Get("x")
		return !ok
	}, time.Second, 50*time.Millisecond)
	require.Equal(t, 0, len(m.LoadAll()))
}

func TestDelete_RemovesKey(t *testing.T) {
	m := utils.NewTTLMap[string, int](time.Minute, time.Minute)
	defer m.Close()

	m.Put("k", 1)
	require.Equal(t, 1, m.Len())

	m.Delete("k")
	_, ok := m.Get("k")
	require.False(t, ok)
	require.Equal(t, 0, m.Len())
}

func TestBackgroundCleanup_RemovesExpired_WithoutAccess(t *testing.T) {
	maxTTL := 150 * time.Millisecond
	cleanup := 50 * time.Millisecond
	m := utils.NewTTLMap[string, int](maxTTL, cleanup)
	defer m.Close()

	// put multiple keys
	for i := 0; i < 10; i++ {
		m.Put("k"+strconv.Itoa(i), i)
	}
	require.Equal(t, 10, m.Len())

	// wait beyond TTL + a couple cleanup ticks; janitor should prune them
	time.Sleep(maxTTL + 3*cleanup)
	require.Equal(t, 0, m.Len())
}

func TestGet_PrunesExpired_OnAccess(t *testing.T) {
	maxTTL := 100 * time.Millisecond
	cleanup := time.Hour // janitor effectively idle
	m := utils.NewTTLMap[string, int](maxTTL, cleanup)
	defer m.Close()

	// Make "old" older than "alive"
	m.Put("old", 2)
	time.Sleep(maxTTL / 2) // age "old" by half the TTL
	m.Put("alive", 1)

	// wait so that "old" expires but "alive" does not
	time.Sleep(maxTTL/2 + 10*time.Millisecond)

	// touching "old" should prune it
	if _, ok := m.Get("old"); ok {
		require.Fail(t, "expected old to be expired")
	}

	// "alive" should still be alive; refresh it
	if _, ok := m.GetAndRefresh("alive"); !ok {
		require.Fail(t, "expected alive to be present")
	}

	require.Equal(t, 1, m.Len())
}

func TestClose_IsIdempotent_AndStopsJanitor(t *testing.T) {
	m := utils.NewTTLMap[string, int](time.Millisecond*10, time.Millisecond*5)
	// Multiple closes should not panic or block forever
	m.Close()
	m.Close()
}

func TestWithContext_CancelStopsJanitor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	m := utils.NewTTLMapWithContext[string, int](ctx, 100*time.Millisecond, 10*time.Millisecond)

	// cancel should allow worker to exit, close should return quickly
	cancel()
	m.Close()
}

func TestDo_DoesNotHoldLock_DuringUserFunction(t *testing.T) {
	m := utils.NewTTLMap[string, int](time.Minute, time.Minute)
	defer m.Close()

	m.Put("a", 1)

	done := make(chan struct{}, 1)
	ok := m.Do("a", func(v int) {
		// re-enter map inside Do callback; should not deadlock
		m.Put("b", v+1)
		close(done)
	})
	require.True(t, ok)

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		require.Fail(t, "Do callback appears to deadlock (lock held during fn?)")
	}

	val, ok := m.Get("b")
	require.True(t, ok)
	require.Equal(t, 2, val)
}
