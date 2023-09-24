package main

import (
	"fmt"

	pokedex_api "github.com/nicwilliams1/pokedexcli/internal"
)

func commandPokedex(c *config, cache *pokedex_api.Cache, s *GameState) error {
	pokedex := s.user.Pokedex
	fmt.Printf("You have %v pokemon in your deck:\n", len(pokedex))
	for _, pokemon := range pokedex {
		fmt.Printf(" -%s\n", pokemon.Name)
	}
	return nil
}
