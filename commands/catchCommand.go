package commands

import (
	"errors"
	"fmt"

	"math/rand"

	"github.com/timeskeletor/pokedexcli/config"
)

type catchCommand struct {}

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
	pkmn, err := cfg.PokeapiClient.GetPokemonSpecies(name)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pkmn.Name)

	roll := rand.Intn(256)
	isShiny := rand.Intn(4096) == 0

	if roll <= pkmn.CaptureRate {
		fmt.Println("Gotcha! You caught", pkmn.Name)
		if isShiny {
			fmt.Println("✨ Whoa! It's a shiny", pkmn.Name + "!")
		}
	
		cfg.CaughtPkmn[pkmn.Name] = pkmn
		cfg.PokeDb.RegisterPokemon(pkmn, true, isShiny)
		return nil
	
	} else {
		cfg.PokeDb.RegisterPokemon(pkmn, false, isShiny)
		fmt.Println(pkmn.Name, "escaped!")
	}

		
	return nil
}