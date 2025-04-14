package utils_test

import (
	"go_project_template/internal/utils"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTTLMap(t *testing.T) {
	t.Run("basic put and get operations", func(t *testing.T) {
		maxTTL := 1 * time.Second
		cleanupInterval := 1 * time.Second
		m := utils.NewTTLMap[int](maxTTL, cleanupInterval)

		// Put a value
		m.Put("key", 42)
		require.Equal(t, 1, m.Len(), "expected TTLMap length to be 1")

		// Get the value
		val, exists := m.Get("key")
		require.True(t, exists, "expected value to exist for key 'key', but it doesn't")
		require.Equal(t, 42, val, "expected value 42, got %d", val)

		// Wait for the cleanup routine to run
		// Check if the value is still present after cleanup
		require.Eventually(t, func() bool {
			_, exists = m.Get("key")
			return !exists
		}, 5*cleanupInterval, cleanupInterval, "expected value to be deleted after TTL expiration, but it still exists")
	})

	t.Run("concurrent put and get operations", func(t *testing.T) {
		maxTTL := 2 * time.Second
		cleanupInterval := 1 * time.Second
		m := utils.NewTTLMap[int](maxTTL, cleanupInterval)

		// Use a wait group to wait for goroutines to finish
		var wg sync.WaitGroup

		// Concurrent put operations
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				key := "key" + strconv.Itoa(index)
				m.Put(key, index)
			}(i)
		}

		// Concurrent get operations
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				key := "key" + strconv.Itoa(index)
				val, exists := m.Get(key)
				if exists {
					require.Equal(t, index, val, "expected value %d, got %d", index, val)
				}
			}(i)
		}

		wg.Wait()
	})
}
