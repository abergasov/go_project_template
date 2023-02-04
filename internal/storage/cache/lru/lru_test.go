package lru_test

import (
	"fmt"
	"go_project_template/internal/storage/cache/lru"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	kbBytes = 1024
)

func TestService_L1_LRU_check(t *testing.T) {
	// given
	iterations := 10
	service, err := lru.CreateL1Cache(1)
	require.NoError(t, err)
	data := make([]byte, 100*kbBytes)
	t.Run("not evict popular piece", func(t *testing.T) {
		service.Purge()
		// when
		for i := 0; i < iterations; i++ {
			require.NoError(t, service.Store(fmt.Sprintf("key_%d", i), data))
		}
		// then
		checkHashL1(t, data, service, iterations)
		// simulate like this piece is very popular
		for i := 0; i < 10; i++ {
			res, ok := service.Get(fmt.Sprintf("key_%d", 0))
			require.True(t, ok, "get key_%d", 0)
			require.Equal(t, data, res, "get key_%d", 0)
		}

		require.NoError(t, service.Store("key_abc", data), "cache should evict old item")
		res, ok := service.Get(fmt.Sprintf("key_%d", 0))
		require.True(t, ok, "expect that evicted will be not `key_0`")
		require.Equal(t, data, res, "get key_%d", 0)
	})
	t.Run("evict unpopular piece", func(t *testing.T) {
		service.Purge()
		// when
		for i := 0; i < iterations; i++ {
			require.NoError(t, service.Store(fmt.Sprintf("key_%d", i), data))
		}
		// then
		// simulate like all pieces very popular
		for i := 0; i < 10; i++ {
			checkHashL1(t, data, service, iterations)
		}
		require.NoError(t, service.Store("key_abc", data), "cache should evict old  item")
		res, ok := service.Get(fmt.Sprintf("key_%d", 0))
		require.False(t, ok, "expect that evicted will be not `key_0`")
		require.Nil(t, res, "get key_%d", 0)
	})
}

func TestService_StoreL1(t *testing.T) {
	// check evicted
	t.Run("append piece of same size", func(t *testing.T) {
		service, err := lru.CreateL1Cache(1)
		require.NoError(t, err)
		// when
		data := make([]byte, 100*kbBytes)    // 100kb data
		storeInL1Cache(t, service, 10, data) // put 1000 bytes in cache. Limit is 1MB, so we can put 10 pieces
		// then
		require.NoError(t, service.Store("key_abc", data))
	})
	t.Run("append piece of bigger size", func(t *testing.T) {
		service, err := lru.CreateL1Cache(1)
		require.NoError(t, err)
		// when
		data := make([]byte, 100*kbBytes)    // 100kb data
		storeInL1Cache(t, service, 10, data) // put 1000 bytes in cache. Limit is 1MB, so we can put 10 pieces
		// then
		require.ErrorIs(t, service.Store("key_abc", make([]byte, 100*kbBytes+1)), lru.ErrL1CacheIsFull)
	})
	t.Run("append piece of smaller size", func(t *testing.T) {
		service, err := lru.CreateL1Cache(1)
		require.NoError(t, err)
		// when
		data := make([]byte, 100*kbBytes)    // 100kb data
		storeInL1Cache(t, service, 10, data) // put 1000 bytes in cache. Limit is 1MB, so we can put 10 pieces
		// then
		require.NoError(t, service.Store("key_abc", make([]byte, 100*kbBytes-1)))
	})
}

func storeInL1Cache(t *testing.T, service *lru.Cache, times int, data []byte) {
	for i := 0; i < times; i++ {
		key := fmt.Sprintf("key_%d", i)
		require.NoErrorf(t, service.Store(key, data), "unexpected error for `%s`", key)
	}
}

func checkHashL1(t *testing.T, data []byte, service *lru.Cache, num int) {
	for i := 0; i < num; i++ {
		res, ok := service.Get(fmt.Sprintf("key_%d", i))
		require.True(t, ok, "get key_%d", i)
		require.Equal(t, data, res, "get key_%d", i)
	}
}

func BenchmarkService_GetPieceFromL1_WithLRU(b *testing.B) {
	service, err := lru.CreateL1Cache(1)
	require.NoError(b, err)
	for i := 0; i < b.N; i++ {
		if err := service.Store(fmt.Sprintf("key%d", i), []byte("123456")); err != nil {
			b.Fatal(err)
		}
	}
	for i := 0; i < b.N; i++ {
		data, ok := service.Get(fmt.Sprintf("key%d", i))
		if !ok {
			b.Fatal("not found")
		}
		if data == nil {
			b.Fatal("data is nil")
		}
	}
}

func TestL1Cache_Resize(t *testing.T) {
	service, err := lru.CreateL1Cache(1)
	require.NoError(t, err)
	payload := make([]byte, 100*kbBytes) // 100kb data
	for i := 0; i < 10; i++ {
		require.NoError(t, service.Store(fmt.Sprintf("key%d", i), payload))
	}

	t.Run("resize to bigger", func(t *testing.T) {
		require.NoError(t, service.Resize(2))
		for i := 0; i < 10; i++ {
			_, ok := service.Get(fmt.Sprintf("key%d", i))
			require.True(t, ok)
		}
		for i := 10; i < 20; i++ {
			require.NoError(t, service.Store(fmt.Sprintf("key%d", i), payload))
		}
		for i := 0; i < 20; i++ {
			_, ok := service.Get(fmt.Sprintf("key%d", i))
			require.True(t, ok)
		}
	})
	t.Run("resize to smaller", func(t *testing.T) {
		require.NoError(t, service.Resize(1))
		for i := 0; i < 10; i++ {
			_, ok := service.Get(fmt.Sprintf("key%d", i))
			require.False(t, ok)
		}
		for i := 11; i < 20; i++ {
			_, ok := service.Get(fmt.Sprintf("key%d", i))
			require.True(t, ok)
		}
	})
}
