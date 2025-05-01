package pokeapi

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache: make(map[string]cacheEntry),
	} 

	go c.reapLoop(interval)
	return c
}

func(c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	newEntry := cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
	c.cache[key] = newEntry
}

func(c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return value.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			<-ticker.C

			now := time.Now()

			c.mu.Lock()

			for key, entry := range c.cache {
				if now.Sub(entry.createdAt) > interval {
					delete(c.cache, key)
				}
			}

			c.mu.Unlock()
		}
	}()
}