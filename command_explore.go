package main

import (
	"fmt"

	pokedex_api "github.com/nicwilliams1/pokedexcli/internal"
)

func commandExplore(c *config, cache *pokedex_api.Cache, s *GameState) error {

	fmt.Printf("Exploring %s...\n", s.area.name)
	url := fmt.Sprintf("%s/%s", c.baseurl, s.area.name)

	_, exists := cache.Get(url)
	if !exists {
		newBody, err := pokedex_api.GetApiResponse(url)
		if err != nil {
			return err
		}
		done := make(chan bool)
		go cache.Add(url, newBody, done)
		<-done
	}

	body, _ := cache.Get(url)
	pokemon, err := pokedex_api.GetPokemonInArea(body)

	if err != nil {
		return err
	}

	fmt.Printf("Found %v pokemon:\n", len(pokemon))
	for _, pokemonName := range pokemon {
		fmt.Printf(" - %s\n", pokemonName)
	}

	s.area.SetCurrentArea(s.area.name, pokemon)

	fmt.Println("")

	return nil
}
