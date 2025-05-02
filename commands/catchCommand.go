package commands

import (
	"fmt"
	"errors"

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
	pkmn, err := cfg.PokeapiClient.GetPokemon(name)
	if err != nil {
		return err
	}
	fmt.Println("Throwing a Pokeball at", pkmn.Name, "...")


	return nil
}