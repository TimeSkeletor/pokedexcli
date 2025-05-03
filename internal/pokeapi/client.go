package pokeapi

import (
	"net/http"
	"time"

	"github.com/timeskeletor/pokedexcli/internal/pokecache"
)

// Client -
type Client struct {
	cache      		*pokecache.Cache[[]byte]
	imageCache      *pokecache.Cache[[]string]
	caughtPkmn      *pokecache.Cache[[]Pokemon]
	httpClient 		http.Client
}

// NewClient -
func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		cache: pokecache.NewCache[[]byte](cacheInterval),
		imageCache: pokecache.NewCache[[]string](cacheInterval),
		caughtPkmn: pokecache.NewCache[[]Pokemon](cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}