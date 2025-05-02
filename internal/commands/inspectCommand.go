package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/timeskeletor/pokedexcli/config"
	"github.com/timeskeletor/pokedexcli/internal/pokeapi"
)

type inspectCommand struct{}

func (c inspectCommand) getName() string {
    return "inspect"
}

func (c inspectCommand) getDescription() string {
    return "Inspect a Pokemon after catching it"
}

func (c inspectCommand) GetCallback(cfg *config.Config, args ...string) error {
    if len(args) != 1 {
        return errors.New("you must provide a Pokemon")
    }

    name := args[0]
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

	var pkmn pokeapi.Pokemon

	pokemon, ok := cfg.CaughtPkmn[name]
	if ok {
		pkmn = pokemon
	} else {
		pokemon, err := cfg.PokeapiClient.GetPokemon(name)
		if err != nil {
			return fmt.Errorf("failed to fetch Pokémon %s: %w", name, err)
		}
		pkmn = pokemon
	}


	caught := cfg.PokeDb.FetchPokemon(ctx, "pokemon", pkmn.ID)

	switch caught {
	case 0:
		fmt.Println("You did not catch a ", pkmn.Name+"yet.")
		return nil
	case 1:
		cfg.CaughtPkmn[pkmn.Name] = pkmn
		fmt.Println("Name:", pkmn.Name)
		fmt.Println("Height:", pkmn.Height)
		fmt.Println("Weight:", pkmn.Weight)
		fmt.Println("Stats:")
		for _, stat := range pkmn.Stats {
			fmt.Printf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, typeInfo := range pkmn.Types {
			fmt.Println("  -", typeInfo.Type.Name)
		}
		return nil
	}

	return nil
}