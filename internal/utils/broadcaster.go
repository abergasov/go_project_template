package utils

import "sync"

type Broadcaster[T any] struct {
	mu              sync.RWMutex
	messages        map[string]chan T
	activeListeners uint64
}

func NewBroadcaster[T any]() *Broadcaster[T] {
	return &Broadcaster[T]{
		messages: make(map[string]chan T),
	}
}

func (b *Broadcaster[T]) RegisterListener(key string) chan T {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.messages[key] = make(chan T)
	b.activeListeners++
	return b.messages[key]
}

func (b *Broadcaster[T]) UnregisterListener(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	close(b.messages[key])
	delete(b.messages, key)
	b.activeListeners--
}

func (b *Broadcaster[T]) Broadcast(msg T) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for i := range b.messages {
		b.messages[i] <- msg
	}
}
