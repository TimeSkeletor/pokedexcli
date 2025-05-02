package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/timeskeletor/pokedexcli/config"
	"github.com/timeskeletor/pokedexcli/commands"
)

func startRepl(cfg *config.Config) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}

		command, exists := commands.GetCommands()[commandName]
		if exists {
			err := command.GetCallback(cfg, args...)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
