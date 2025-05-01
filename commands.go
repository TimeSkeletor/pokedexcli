package main

import (
	"fmt"
	"os"
	"errors"
)

func getCommands() map[string]cliCommand {
    return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Get instructions to use the Pokedex",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Get the next page of locations",
			callback:    commandMapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous page of locations",
			callback:    commandMapb,
		},
        "exit": {
            name:        "exit",
            description: "Exit the Pokedex",
            callback:    commandExit,
        },
    }
}


func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!\n")
	fmt.Println("Usage:\n")

	for _, cmd := range getCommands() {
		fmt.Printf("- %s - %s\n", cmd.name, cmd.description) 
	}
	fmt.Println()
	return nil
}


func commandMapf(cfg *config) error {
	locationsResp, err := cfg.pokeapiClient.GenericGetList(cfg.nextLocationsURL, "location-area")
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locationsResp.Next
	cfg.prevLocationsURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(cfg *config) error {
	if cfg.prevLocationsURL == nil {
		return errors.New("you're on the first page")
	}

	locationResp, err := cfg.pokeapiClient.GenericGetList(cfg.prevLocationsURL, "location-area")
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locationResp.Next
	cfg.prevLocationsURL = locationResp.Previous

	for _, loc := range locationResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandExit(cfg *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandNotFound(cfg *config) error {
	return errors.New("command not found")
}
