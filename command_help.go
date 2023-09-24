package main

import (
	"fmt"

	pokedex_api "github.com/nicwilliams1/pokedexcli/internal"
)

func commandHelp(c *config, cache *pokedex_api.Cache, s *GameState) error {

	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()

	return nil
}
