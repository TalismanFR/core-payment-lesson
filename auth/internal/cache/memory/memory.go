package memory

import (
	"auth/internal/cache"
	"auth/internal/model/principal"
	"errors"
	"sync"
	"time"
)

var (
	ErrNoKey = errors.New("key wasn't found")
)

var _ cache.Cache = (*InmemoryCache)(nil)

type InmemoryCache struct {
	m          sync.Map
	lock       sync.RWMutex
	entryCount int64
	hitCount   int64
	missCount  int64
	expireIn   time.Duration
}

func New(expireIn time.Duration) *InmemoryCache {
	return &InmemoryCache{expireIn: expireIn}
}

func (i *InmemoryCache) StartCleaning(period time.Duration) (cancel func()) {
	done := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-time.After(period):
				i.m.Range(func(key, value any) bool {
					e, _ := value.(entry)
					if e.expired() {

						// We would like to delete expired entry,
						// but we don't want to wait
						go i.del(key.(string))
					}
					return true
				})
			}
		}
	}()

	return func() {
		done <- struct{}{}
	}
}

func (i *InmemoryCache) Get(key string) (p principal.Principal, err error) {

	defer func() {
		i.hitOrMiss(err == nil)
	}()

	v, ok := i.m.Load(key)
	if !ok {
		return p, ErrNoKey
	}

	e, _ := v.(entry)
	if e.expired() {

		// We would like to delete expired entry,
		// but we don't want to wait
		go i.del(key)

		return p, ErrNoKey
	}

	return e.p, nil
}

func (i *InmemoryCache) Set(key string, p principal.Principal) error {

	_, loaded := i.m.Load(key)
	if !loaded {
		i.lock.Lock()
		i.entryCount++
		i.lock.Unlock()
	}

	i.m.Store(key, entry{
		p:        p,
		expireAt: time.Now().Add(i.expireIn),
	})

	return nil
}

func (i *InmemoryCache) Del(key string) bool {
	return i.del(key)
}

func (i *InmemoryCache) EntryCount() int64 {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.entryCount
}

func (i *InmemoryCache) HitCount() int64 {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.hitCount
}

func (i *InmemoryCache) MissCount() int64 {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.missCount
}

func (i *InmemoryCache) hitOrMiss(hit bool) {
	i.lock.Lock()
	defer i.lock.Unlock()
	if hit {
		i.hitCount++
	} else {
		i.missCount++
	}
}

func (i *InmemoryCache) del(key string) bool {

	_, affected := i.m.LoadAndDelete(key)

	if affected {
		i.lock.Lock()
		i.entryCount--
		i.lock.Unlock()
	}

	return affected
}

type entry struct {
	p        principal.Principal
	expireAt time.Time
}

func (e *entry) expired() bool {
	return time.Since(e.expireAt) >= 0
}
