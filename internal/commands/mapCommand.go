package commands

import (
	"fmt"
	"errors"

	"github.com/timeskeletor/pokedexcli/config"
)

type mapfCommand struct {}

func (c mapfCommand) getName() string {
	return "map"
}

func (c mapfCommand) getDescription() string {
	return "Get the next page of locations"
}

func (c mapfCommand) GetCallback(cfg *config.Config, args ...string) error {
	locationsResp, err := cfg.PokeapiClient.GenericGetList(cfg.NextPageURL, "location-area")
	if err != nil {
		return err
	}

	cfg.NextPageURL = locationsResp.Next
	cfg.PrevPageURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

type mapbCommand struct {}


func (c mapbCommand) getName() string {
	return "mapb"
}

func (c mapbCommand) getDescription() string {
	return "Get the previous page of locations"
}

func (c mapbCommand) GetCallback(cfg *config.Config, args ...string) error {
	if cfg.PrevPageURL == nil {
		return errors.New("you're on the first page")
	}

	locationResp, err := cfg.PokeapiClient.GenericGetList(cfg.PrevPageURL, "location-area")
	if err != nil {
		return err
	}

	cfg.NextPageURL = locationResp.Next
	cfg.PrevPageURL = locationResp.Previous

	for _, loc := range locationResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}
