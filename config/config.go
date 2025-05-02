package config

import (
	"github.com/timeskeletor/pokedexcli/internal/pokeapi"
	"github.com/timeskeletor/pokedexcli/internal/pokedatabase"
)


const (
	ReplMsg = "Pokedex > "
)


type Config struct {
	PokeapiClient	 pokeapi.Client
	PokeDb			*pokedatabase.Database
	CaughtPkmn		map[string]pokeapi.PokemonSpecies
	NextPageURL 	*string
	PrevPageURL 	*string
}
