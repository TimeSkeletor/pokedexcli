package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/timeskeletor/pokedexcli/config"
)

type pokedexCommand struct{}

func (c pokedexCommand) getName() string {
    return "pokedex"
}

func (c pokedexCommand) getDescription() string {
    return "Lists all the pokemon you caught"
}

func (c pokedexCommand) GetCallback(cfg *config.Config, args ...string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

	results, _ := cfg.PokeDb.FetchAllRegisteredPokemon(ctx, "pokemon")
	fmt.Println("Your Pokedex:")

	for _, result := range results {
		fmt.Println("- "+result)
	}

    return nil
}