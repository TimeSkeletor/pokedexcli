package commands

import (
	"fmt"
	"errors"
	
	"github.com/timeskeletor/pokedexcli/config"
)

type exploreCommand struct {}

func (c exploreCommand) getName() string {
	return "explore <location_name>"
}

func (c exploreCommand) getDescription() string {
	return "Explore a location"
}

func (c exploreCommand) GetCallback(cfg *config.Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a location name")
	}

	name := args[0]
	location, err := cfg.PokeapiClient.GetLocation(name)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", location.Name)
	fmt.Println("Found Pokemon: ")
	for _, enc := range location.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}
	return nil
}