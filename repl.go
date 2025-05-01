package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/timeskeletor/pokedexcli/internal/pokeapi"
)

const (
	replMsg = "Pokedex > "
)
type config struct {
	pokeapiClient    pokeapi.Client
	nextPageURL *string
	prevPageURL *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}


func startRepl(cfg *config) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(replMsg)
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]

		command, exists := getCommands()[commandName]
		if exists {
			err := command.callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println(replMsg, "Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
