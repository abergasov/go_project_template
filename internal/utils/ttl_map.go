package utils

import (
	"context"
	"hash/maphash"
	"sync"
	"time"
)

// numShards controls lock striping granularity; must be a power of two.
const (
	numShards = 32
	shardMask = numShards - 1
)

// item stores the payload alongside its expiry as a monotonic int64 nanosecond
// timestamp. Using int64 instead of time.Time halves the field size (8 vs 24 bytes)
// and reduces comparison cost to a single integer compare.
type item[V any] struct {
	value      V
	expiryNano int64
}

type cacheShard[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]item[V]
	// this extra field will guarantee that object will be 64 bytes size, tweak for efficient cache usage
	_ [32]byte // pads to 64 B (1 cache line) so adjacent shards don't cause false sharing
}

// TTLMap stores values with a time-to-live.
// Expired items are removed during cleanup or upon access.
//
// Performance characteristics vs a single-mutex map:
//   - 32 independent shards → up to 32× lower lock contention on concurrent reads/writes
//   - Value-typed items (no pointer per entry) → lower GC overhead and better locality
//   - int64 nanosecond expiry → cheaper comparison, 8 bytes instead of 24
//   - Get uses RLock on the hot (non-expired) path → concurrent reads scale with cores
type TTLMap[K comparable, V any] struct {
	shards [numShards]cacheShard[K, V]
	seed   maphash.Seed
	maxTTL int64 // nanoseconds

	wg        sync.WaitGroup
	stopCh    chan struct{}
	closeOnce sync.Once
}

// NewTTLMap creates a TTL map with a background cleanup goroutine.
// Call Close() when done to stop the goroutine.
func NewTTLMap[K comparable, V any](maxTTL, cleanupInterval time.Duration) *TTLMap[K, V] {
	return NewTTLMapWithContext[K, V](context.Background(), maxTTL, cleanupInterval)
}

// NewTTLMapWithContext ties the cleanup goroutine to ctx and also supports Close().
func NewTTLMapWithContext[K comparable, V any](ctx context.Context, maxTTL, cleanupInterval time.Duration) *TTLMap[K, V] {
	m := &TTLMap[K, V]{
		seed:   maphash.MakeSeed(),
		maxTTL: maxTTL.Nanoseconds(),
		stopCh: make(chan struct{}),
	}
	for i := range m.shards {
		m.shards[i].m = make(map[K]item[V], 64)
	}

	ticker := time.NewTicker(cleanupInterval)
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				now := time.Now().UnixNano()
				for i := range m.shards {
					s := &m.shards[i]
					s.mu.Lock()
					for k, v := range s.m {
						if now > v.expiryNano {
							delete(s.m, k)
						}
					}
					s.mu.Unlock()
				}
			case <-m.stopCh:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
	return m
}

// Close stops the background cleanup goroutine and waits for it to exit, safe to call multiple times.
func (m *TTLMap[K, V]) Close() {
	m.closeOnce.Do(func() {
		close(m.stopCh)
		m.wg.Wait()
	})
}

// shardOf returns the shard responsible for key k.
func (m *TTLMap[K, V]) shardOf(k K) *cacheShard[K, V] {
	return &m.shards[maphash.Comparable(m.seed, k)&shardMask]
}

// Len returns the total number of entries across all shards (including not-yet-evicted expired ones).
// Thread-safe.
func (m *TTLMap[K, V]) Len() int {
	total := 0
	for i := range m.shards {
		s := &m.shards[i]
		s.mu.RLock()
		total += len(s.m)
		s.mu.RUnlock()
	}
	return total
}

// Put stores value under the specified key and refreshes its expiry time.
func (m *TTLMap[K, V]) Put(k K, v V) {
	s := m.shardOf(k)
	s.mu.Lock()
	s.m[k] = item[V]{value: v, expiryNano: time.Now().UnixNano() + m.maxTTL}
	s.mu.Unlock()
}

// Delete removes the value associated with the given key.
func (m *TTLMap[K, V]) Delete(k K) {
	s := m.shardOf(k)
	s.mu.Lock()
	delete(s.m, k)
	s.mu.Unlock()
}

// Get returns the value associated with key.
// Uses RLock on the hot (live) path for concurrent read scalability;
// upgrades to a write lock only when an expired entry must be deleted.
func (m *TTLMap[K, V]) Get(k K) (val V, ok bool) {
	s := m.shardOf(k)

	s.mu.RLock()
	it, ok := s.m[k]
	if !ok {
		s.mu.RUnlock()
		return val, false
	}
	if time.Now().UnixNano() <= it.expiryNano {
		val = it.value
		s.mu.RUnlock()
		return val, true
	}
	s.mu.RUnlock()

	// Expired: upgrade to write lock to evict.
	s.mu.Lock()
	if it2, ok2 := s.m[k]; ok2 && time.Now().UnixNano() > it2.expiryNano {
		delete(s.m, k)
	}
	s.mu.Unlock()
	return val, false
}

// GetAndRefresh returns the value associated with key and resets its TTL.
// If the key is missing or expired, the zero value and false are returned.
func (m *TTLMap[K, V]) GetAndRefresh(k K) (val V, ok bool) {
	s := m.shardOf(k)
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.m[k]
	if !ok || time.Now().UnixNano() > it.expiryNano {
		if ok {
			delete(s.m, k)
		}
		return val, false
	}
	it.expiryNano = time.Now().UnixNano() + m.maxTTL
	s.m[k] = it
	return it.value, true
}

// Do executes fn with the value associated with key.
// The map lock is released before fn is called, so fn may safely access the map.
// If the key is missing or expired, fn is not called and false is returned.
func (m *TTLMap[K, V]) Do(k K, fn func(v V)) bool {
	s := m.shardOf(k)

	s.mu.RLock()
	it, ok := s.m[k]
	s.mu.RUnlock()
	if !ok || time.Now().UnixNano() > it.expiryNano {
		if ok {
			s.mu.Lock()
			if it2, ok2 := s.m[k]; ok2 && time.Now().UnixNano() > it2.expiryNano {
				delete(s.m, k)
			}
			s.mu.Unlock()
		}
		return false
	}
	fn(it.value)
	return true
}

// DoAndApply modifies the value associated with key using fn.
// The expiry time is not changed. Returns false if the key is missing or expired.
func (m *TTLMap[K, V]) DoAndApply(k K, fn func(v V) V) bool {
	s := m.shardOf(k)
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.m[k]
	if !ok {
		return false
	}
	if time.Now().UnixNano() > it.expiryNano {
		delete(s.m, k)
		return false
	}
	it.value = fn(it.value)
	s.m[k] = it
	return true
}

// Upsert inserts zeroValue if the key does not exist, or updates the existing value
// using fn. Either way the TTL is refreshed.
func (m *TTLMap[K, V]) Upsert(key K, fn func(value V) V, zeroValue V) {
	s := m.shardOf(key)
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.m[key]
	if !ok {
		s.m[key] = item[V]{value: zeroValue, expiryNano: time.Now().UnixNano() + m.maxTTL}
		return
	}
	s.m[key] = item[V]{value: fn(it.value), expiryNano: time.Now().UnixNano() + m.maxTTL}
}

// LoadAll returns a snapshot of all live (non-expired) entries.
func (m *TTLMap[K, V]) LoadAll() map[K]V {
	now := time.Now().UnixNano()
	result := make(map[K]V)
	for i := range m.shards {
		s := &m.shards[i]
		s.mu.RLock()
		for k, v := range s.m {
			if now <= v.expiryNano {
				result[k] = v.value
			}
		}
		s.mu.RUnlock()
	}
	return result
}
