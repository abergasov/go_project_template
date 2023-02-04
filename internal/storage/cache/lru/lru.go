package lru

import (
	"errors"
	"fmt"
	"go_project_template/internal/storage/cache/lru/utils"

	lru "github.com/hashicorp/golang-lru"
)

const (
	bytesInMB = 1024 * 1024 // equal 1_048_576
)

var (
	ErrL1CacheIsFull = errors.New("l1 cache is full")
)

type Cache struct {
	lruCache       *lru.Cache        // LRU cache only satisfies for cdn
	sizerL1        *utils.CacheSizer // LRU library does not track size of items, so we need to track it manually
	evictedItems   chan EvictedItem
	maxSizeInBytes int64
}

func CreateL1Cache(lruMaxSizeMB int64) (*Cache, error) {
	// Get MaxSize in MB and convert it to bytes.
	// This is max size of L1 cache in items (as we suggest, that each item from 1 byte to 100MB)
	// If simply memory divine to 100MB - that it will not be efficient, because we can have a lot of small items.
	// So max size of L1 cache in items is 100MB * MaxSize.
	// l1BytesSize is used to track size of items in L1 cache and prevent memory overflow.
	maxSizeInBytes := lruMaxSizeMB * bytesInMB
	l1 := &Cache{
		sizerL1:        utils.CreateCacheSizer(maxSizeInBytes),
		maxSizeInBytes: maxSizeInBytes,
	}
	lruCache, err := lru.NewWithEvict(int(maxSizeInBytes), l1.onEvict)
	if err != nil {
		return nil, fmt.Errorf("error while creating LRU cache: %w", err)
	}
	l1.lruCache = lruCache
	return l1, nil
}

func (s *Cache) GetMaxSizeBytes() int64 {
	return s.maxSizeInBytes
}

func (s *Cache) EvictListener() <-chan EvictedItem {
	if s.evictedItems == nil {
		s.evictedItems = make(chan EvictedItem, 1000)
	}
	return s.evictedItems
}

// Store stores data in L1 cache. If L1 cache is full, returns ErrL1CacheIsFull.
func (s *Cache) Store(key string, value []byte) error {
	// check that storage is less than maximum size
	if s.sizerL1.CanAddWithoutEvicting(len(value)) {
		// we are lucky, just add it to object storage
		s.lruCache.Add(key, value)
		s.sizerL1.Add(len(value))
		return nil
	}
	// cache will evict some items
	_, valueOld, ok := s.lruCache.GetOldest()
	if !ok {
		// there are no items to evict, so adding this item will increase size of storage
		return ErrL1CacheIsFull
	}
	valueOldBytes, ok := valueOld.([]byte)
	if !ok {
		// unknown type of value, skip it
		return ErrL1CacheIsFull
	}
	if len(valueOldBytes) < len(value) {
		// new item is bigger than old item, so we can't replace it
		return ErrL1CacheIsFull
	}
	// we can replace old item with new item
	s.lruCache.RemoveOldest() // as we manually track size of items, we remove onEvict
	s.lruCache.Add(key, value)
	s.sizerL1.Add(len(value))
	return nil
}

// Purge removes all items from L1 cache and reset bytes counter.
func (s *Cache) Purge() {
	s.lruCache.Purge()
	s.sizerL1.Purge()
}

func (s *Cache) Get(key string) (value []byte, exist bool) {
	if data, ok := s.lruCache.Get(key); ok {
		if value, ok = data.([]byte); ok {
			return value, true
		}
	}
	return nil, false
}

func (s *Cache) onEvict(key, value interface{}) {
	valueBytes, ok := value.([]byte)
	if !ok {
		// unknown type of value, skip it
		return
	}
	s.sizerL1.Remove(len(valueBytes))
	if s.evictedItems != nil {
		s.evictedItems <- EvictedItem{Key: key.(string), Value: valueBytes}
	}
}

func (s *Cache) Resize(sizeMB int64) error {
	newSizeInBytes := sizeMB * bytesInMB
	if s.maxSizeInBytes < newSizeInBytes {
		// we simply increase size of L2 cache without any eviction
		s.maxSizeInBytes = newSizeInBytes
		s.sizerL1.Resize(int(newSizeInBytes))
		return nil
	}
	// we need to evict some items from L2 cache
	delta := s.maxSizeInBytes - newSizeInBytes
	evictKeysSize := 0
	for {
		_, value, ok := s.lruCache.GetOldest()
		if !ok {
			// there are no items to evict, so adding this item will increase size of storage
			break
		}
		valueBytes, ok := value.([]byte)
		if !ok {
			// we expect only []byte, so trigger error here
			return errors.New("unknown type of value in L2 cache")
		}
		s.lruCache.RemoveOldest() // it will call onEvict
		s.sizerL1.Remove(len(valueBytes))
		evictKeysSize += len(valueBytes)
		if evictKeysSize >= int(delta) {
			break
		}
	}
	s.maxSizeInBytes = newSizeInBytes
	s.lruCache.Resize(int(s.maxSizeInBytes))
	return nil
}
