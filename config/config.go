package config

import (
	"github.com/timeskeletor/pokedexcli/internal/pokeapi"
)


const (
	ReplMsg = "Pokedex > "
)


type Config struct {
	PokeapiClient    pokeapi.Client
	NextPageURL *string
	PrevPageURL *string
}
