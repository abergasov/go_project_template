package utils

import "sync"

type RWSlice[V any] struct {
	mu       sync.RWMutex
	Internal []V
}

func NewRWSlice[V any]() *RWSlice[V] {
	return &RWSlice[V]{
		Internal: make([]V, 0),
	}
}

func (rw *RWSlice[V]) Add(internal V) {
	rw.mu.Lock()
	rw.Internal = append(rw.Internal, internal)
	rw.mu.Unlock()
}

func (rw *RWSlice[V]) AddBulk(internal []V) {
	rw.mu.Lock()
	rw.Internal = append(rw.Internal, internal...)
	rw.mu.Unlock()
}

func (rw *RWSlice[V]) LoadAll() []V {
	rw.mu.RLock()
	actual := rw.Internal
	rw.mu.RUnlock()
	return actual
}

func (rw *RWSlice[V]) LoadAndErase() []V {
	rw.mu.Lock()
	actual := rw.Internal
	rw.Internal = make([]V, 0, len(actual))
	rw.mu.Unlock()
	return actual
}

func (rw *RWSlice[V]) Len() int {
	rw.mu.RLock()
	defer rw.mu.RUnlock()
	return len(rw.Internal)
}

func (rw *RWSlice[V]) Erase() {
	rw.mu.Lock()
	rw.Internal = make([]V, 0, len(rw.Internal))
	rw.mu.Unlock()
}
