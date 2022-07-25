package cache

//type Cache interface {
//
//	// Get returns the value or not found error.
//	Get(key []byte) ([]byte, error)
//
//	// Set sets a key, value and expiration for a cache entry and stores it in the cache.
//	// expireIn <= 0 means no expire, but it can be evicted when cache is full.
//	Set(key []byte, val []byte, expireIn int64) error
//
//	// Del deletes an item in the cache by key and returns true or false if a deletion occurred.
//	Del(key []byte) (affected bool)
//
//	// EntryCount returns the number of items currently in the cache.
//	EntryCount() (entryCount int64)
//	// HitCount is a metric that returns number of times a key was found in the cache.
//	HitCount() int64
//	// MissCount is a metric that returns the number of times a miss occurred in the cache.
//	MissCount() int64
//}
