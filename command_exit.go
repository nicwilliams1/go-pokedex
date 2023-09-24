package main

import (
	"os"

	pokedex_api "github.com/nicwilliams1/pokedexcli/internal"
)

func commandExit(c *config, cache *pokedex_api.Cache, s *GameState) error {
	os.Exit(0)
	return nil
}
