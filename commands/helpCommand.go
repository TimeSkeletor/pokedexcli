package commands

import (
	"fmt"
	
	"github.com/timeskeletor/pokedexcli/config"
)

type helpCommand struct {}

func (c helpCommand) getName() string {
	return "help"
}

func (c helpCommand) getDescription() string {
	return "Displays a help message"
}

func (c helpCommand) GetCallback(cfg *config.Config, args ...string) error {
	fmt.Println("Usage:")

	for _, cmd := range GetCommands() {
		fmt.Printf("- %s - %s\n", cmd.getName(), cmd.getDescription()) 
	}
	fmt.Println()
	return nil
}