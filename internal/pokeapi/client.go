package pokeapi

import (
	"net/http"
	"time"

	"github.com/timeskeletor/pokedexcli/internal/pokeapi"
)

// Client -
type Client struct {
	httpClient http.Client
	cache *pokeapi.Cache
}

// NewClient -
func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokeapi.NewCache(5 * time.Minute),
	}
}