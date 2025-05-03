package commands

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/timeskeletor/pokedexcli/config"
)

type catchCommand struct{}

func (c catchCommand) getName() string {
    return "catch"
}

func (c catchCommand) getDescription() string {
    return "Catch a Pokemon"
}

func (c catchCommand) GetCallback(cfg *config.Config, args ...string) error {
    if len(args) != 1 {
        return errors.New("you must provide a Pokemon")
    }

    name := args[0]
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    pkmn, err := cfg.PokeapiClient.GetPokemonSpecies(name)
    if err != nil {
        return fmt.Errorf("failed to fetch Pokémon %s: %w", name, err)
    }

    // Check if Pokémon is already registered
    caught := cfg.PokeDb.FetchPokemon(ctx, "pokemon", pkmn.ID)
    fmt.Printf("A wild %s, appeared! Is caught: %d\n", pkmn.Name, caught)
    cfg.PokeDb.RegisterPokemon(ctx, pkmn)

    switch caught {
    case 1:
        fmt.Printf("%s is already caught!\n", pkmn.Name)
        return nil
    case 0:
        fmt.Printf("Throwing a Pokeball at %s...\n", pkmn.Name)
        roll := rand.Intn(256)
        isShiny := rand.Intn(4096) == 0
        captureChance := roll <= pkmn.CaptureRate

        if captureChance {
            fmt.Println("Gotcha! You caught", pkmn.Name)
            cfg.PokeDb.CatchPokemon(ctx, "pokemon", pkmn.ID, isShiny)
            if isShiny {
                fmt.Println("✨ Whoa! It's a shiny", pkmn.Name+"!")
            }
            fmt.Println("You may now inspect it with the inspect command.")
            return nil
        }
        fmt.Println(pkmn.Name, "escaped!")
        return nil
    }

    return nil
}