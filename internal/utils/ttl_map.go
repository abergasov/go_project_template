package utils

import (
	"sync"
	"time"
)

type item[T any] struct {
	value      T
	expiryTime time.Time
}

// TTLMap map with items that will be deleted after maxTTL seconds
type TTLMap[T any] struct {
	mu     sync.RWMutex
	m      map[string]*item[T]
	maxTTL time.Duration
}

// NewTTLMap creates new map with items that will be deleted after maxTTL seconds
func NewTTLMap[T any](maxTTL, cleanupInterval time.Duration) (m *TTLMap[T]) {
	m = &TTLMap[T]{
		m:      make(map[string]*item[T], 1_000),
		maxTTL: maxTTL,
	}
	ticker := time.NewTicker(cleanupInterval)
	go func() {
		defer ticker.Stop()
		for now := range ticker.C {
			m.mu.Lock()
			for k, v := range m.m {
				if now.After(v.expiryTime) {
					delete(m.m, k)
				}
			}
			m.mu.Unlock()
		}
	}()
	return
}

// Len returns the number of items in the map
func (m *TTLMap[T]) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.m)
}

// Put sets the given value under the specified key.
func (m *TTLMap[T]) Put(k string, v T) {
	m.mu.Lock()
	defer m.mu.Unlock()

	it, ok := m.m[k]
	if !ok {
		it = &item[T]{value: v}
		m.m[k] = it
	} else {
		it.value = v
	}
	it.expiryTime = time.Now().Add(m.maxTTL)
}

// Delete removes the value associated with the given key.
func (m *TTLMap[T]) Delete(k string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, k)
}

// Get gets the value associated with the given key.
func (m *TTLMap[T]) Get(k string) (val T, exists bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if it, ok := m.m[k]; ok {
		return it.value, ok
	}
	return val, false
}

// GetAndRefresh gets the value associated with the given key and refreshes its expiry time.
func (m *TTLMap[T]) GetAndRefresh(k string) (val T, exists bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if it, ok := m.m[k]; ok {
		it.expiryTime = time.Now().Add(m.maxTTL)
		return it.value, true
	}
	return val, false
}
