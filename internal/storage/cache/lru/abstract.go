package lru

type EvictedItem struct {
	Key   string
	Value []byte
}

type Cacher interface {
	Store(key string, value []byte) error
	Get(key string) (value []byte, exist bool)
	EvictListener() <-chan EvictedItem
	Purge()
	Resize(sizeMB int64) error
	GetMaxSizeBytes() int64
}
