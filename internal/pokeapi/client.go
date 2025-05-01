package pokeapi

import (
	"net/http"
	"time"
)

// Client -
type Client struct {
	httpClient http.Client
	cache *Cache
}

// NewClient -
func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: NewCache(5 * time.Minute),
	}
}