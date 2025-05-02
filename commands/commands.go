package commands

import (
	"github.com/timeskeletor/pokedexcli/config"
)


type Command interface {
	getName() string
	getDescription() string
	GetCallback(cfg *config.Config, args ...string) error
}


func GetCommands() map[string]Command {
	return map[string]Command{
		"catch":   catchCommand{},
		"help":    helpCommand{},
		"explore": exploreCommand{},
		"map": mapfCommand{},
		"mapb": mapbCommand{},
		"exit": exitCommand{},
	}
}
