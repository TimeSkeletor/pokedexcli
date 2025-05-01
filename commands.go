package main

import (
	"fmt"
	"os"
	"errors"
)


const (
	replMsg = "Pokedex > "
)

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Usage:")

	for _, cmd := range getCommands() {
		fmt.Printf("- %s - %s\n", cmd.name, cmd.description) 
	}
	fmt.Println()
	return nil
}


func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a location name")
	}

	name := args[0]
	location, err := cfg.pokeapiClient.GetLocation(name)
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

func commandMapf(cfg *config, args ...string) error {
	locationsResp, err := cfg.pokeapiClient.GenericGetList(cfg.nextPageURL, "location-area")
	if err != nil {
		return err
	}

	cfg.nextPageURL = locationsResp.Next
	cfg.prevPageURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.prevPageURL == nil {
		return errors.New("you're on the first page")
	}

	locationResp, err := cfg.pokeapiClient.GenericGetList(cfg.prevPageURL, "location-area")
	if err != nil {
		return err
	}

	cfg.nextPageURL = locationResp.Next
	cfg.prevPageURL = locationResp.Previous

	for _, loc := range locationResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println(replMsg, "Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

