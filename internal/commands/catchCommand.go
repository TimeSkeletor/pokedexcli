package commands

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"crawshaw.io/sqlite"
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
    _, _, err = cfg.PokeDb.FetchPokemon(ctx, "pokemon", pkmn.ID)
    if err == nil {
        // Pokémon exists, check if it's already caught
        var caught bool
        err = cfg.PokeDb.ExecQuery(ctx, "SELECT caught FROM pokemon WHERE number = ?", func(stmt *sqlite.Stmt) error {
            caught = stmt.GetInt64("caught") == 1
            return nil
        }, pkmn.ID)
        if err != nil {
            return fmt.Errorf("failed to check Pokémon status: %w", err)
        }
        if caught {
            fmt.Printf("%s is already registered!\n", pkmn.Name)
            return nil
        }
        // Update caught status
        fmt.Printf("Throwing a Pokeball at %s...\n", pkmn.Name)
        roll := rand.Intn(256)
        isShiny := rand.Intn(4096) == 0
        if roll <= pkmn.CaptureRate {
            fmt.Println("Gotcha! You caught", pkmn.Name)
            if isShiny {
                fmt.Println("✨ Whoa! It's a shiny", pkmn.Name+"!")
            }
            err = cfg.PokeDb.CatchPokemon(ctx, "pokemon", pkmn.ID)
            if err != nil {
                return fmt.Errorf("failed to update caught status: %w", err)
            }
            cfg.CaughtPkmn[pkmn.Name] = pkmn
            return nil
        }
        fmt.Println(pkmn.Name, "escaped!")
        return nil
    }

    // Pokémon doesn't exist, register it
    fmt.Printf("Throwing a Pokeball at %s...\n", pkmn.Name)
    roll := rand.Intn(256)
    isShiny := rand.Intn(4096) == 0
    caught := roll <= pkmn.CaptureRate
    if caught {
        fmt.Println("Gotcha! You caught", pkmn.Name)
        if isShiny {
            fmt.Println("✨ Whoa! It's a shiny", pkmn.Name+"!")
        }
    } else {
        fmt.Println(pkmn.Name, "escaped!")
    }

    err = cfg.PokeDb.RegisterPokemon(ctx, pkmn, caught, isShiny)
    if err != nil {
        return fmt.Errorf("failed to register Pokémon: %w", err)
    }
    if caught {
        cfg.CaughtPkmn[pkmn.Name] = pkmn
    }
    return nil
}