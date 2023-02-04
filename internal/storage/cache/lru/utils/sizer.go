package utils

import "sync/atomic"

type CacheSizer struct {
	bytesSize   int64 // track manually size of cache
	maxSize     int64 // max size of cache in bytes
	currentSize int64 // current size of cache in items
}

func CreateCacheSizer(maxSizeBytes int64) *CacheSizer {
	return &CacheSizer{
		maxSize: maxSizeBytes,
	}
}

// CanAddWithoutEvicting check that storage is less than maximum size, and we can add item to cache without evicting.
func (c *CacheSizer) CanAddWithoutEvicting(size int) bool {
	return atomic.LoadInt64(&c.bytesSize)+int64(size) < c.maxSize
}

func (c *CacheSizer) Add(size int) {
	atomic.AddInt64(&c.bytesSize, int64(size))
	atomic.AddInt64(&c.currentSize, int64(1))
}

func (c *CacheSizer) Remove(size int) {
	atomic.AddInt64(&c.bytesSize, -int64(size))
	atomic.AddInt64(&c.currentSize, -int64(1))
}

func (c *CacheSizer) Purge() {
	atomic.StoreInt64(&c.bytesSize, 0)
	atomic.StoreInt64(&c.currentSize, 0)
}

func (c *CacheSizer) Resize(bytesSize int) {
	atomic.StoreInt64(&c.maxSize, int64(bytesSize))
}
