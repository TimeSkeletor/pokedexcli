package pokeapi

import (
	"net/http"
	"time"

	"github.com/timeskeletor/pokedexcli/internal/pokecache"
)

// Client -
type Client struct {
	cache      		pokecache.Cache
	caughtPkmn      map[string]Pokemon
	httpClient 		http.Client
}

// NewClient -
func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		cache: pokecache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}