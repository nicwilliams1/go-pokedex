package main

import (
	"errors"
	"fmt"

	pokedex_api "github.com/nicwilliams1/pokedexcli/internal"
)

func commandInspect(c *config, cache *pokedex_api.Cache, s *GameState) error {
	pokemon, ok := s.user.Pokedex[s.inspectTarget]
	if !ok {
		return errors.New("you have not caught that pokemon yet")
	}

	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf("  -%v: %v\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  -%v\n", t.Type.Name)
	}

	return nil
}
