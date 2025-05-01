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

func NewCache(duration int) {
	return 
}