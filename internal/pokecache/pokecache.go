package pokecache

import (
	"sync"
	"time"
)

type Cache[T any] struct {
	cache map[string]cacheEntry[T]
	mux   *sync.Mutex
}

type cacheEntry[T any] struct {
	createdAt time.Time
	val       T
}

func NewCache[T any](interval time.Duration) *Cache[T] {
	c := &Cache[T]{
		cache: make(map[string]cacheEntry[T]),
		mux:   &sync.Mutex{},
	}

	go c.reapLoop(interval)
	return c
}

func (c *Cache[T]) Add(key string, value T) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.cache[key] = cacheEntry[T]{
		createdAt: time.Now(),
		val:       value,
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	entry, ok := c.cache[key]
	return entry.val, ok
}

func (c *Cache[T]) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(time.Now(), interval)
	}
}

func (c *Cache[T]) reap(now time.Time, maxAge time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()

	for k, v := range c.cache {
		if now.Sub(v.createdAt) > maxAge {
			delete(c.cache, k)
		}
	}
}
