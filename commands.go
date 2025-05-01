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
	fmt.Println("Usage:")

	for _, cmd := range getCommands() {
		fmt.Printf("- %s - %s\n", cmd.name, cmd.description) 
	}
	fmt.Println()
	return nil
}


func commandMapf(cfg *config) error {
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

func commandMapb(cfg *config) error {
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

func commandExit(cfg *config) error {
	fmt.Println(replMsg, "Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

