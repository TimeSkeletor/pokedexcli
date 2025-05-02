package main

import (
	"time"

	"github.com/timeskeletor/pokedexcli/config"
	"github.com/timeskeletor/pokedexcli/internal/pokeapi"
	"github.com/timeskeletor/pokedexcli/internal/pokedatabase"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, time.Minute*5)
	cfg := &config.Config{
		PokeapiClient: pokeClient,
	}
	pokedatabase.SetDb()
	startRepl(cfg)
}
