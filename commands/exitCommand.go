package commands

import (
	"fmt"
	"os"

	"github.com/timeskeletor/pokedexcli/config"
)

type exitCommand struct {}

func (c exitCommand) getName() string {
	return "exit"
}

func (c exitCommand) getDescription() string {
	return "Exit the Pokedex"
}

func (c exitCommand) GetCallback(cfg *config.Config, args ...string) error {
	fmt.Println(config.ReplMsg, "Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
