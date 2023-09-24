package main

import (
	pokedex_api "github.com/nicwilliams1/pokedexcli/internal"
)

func commandMapb(c *config, cache *pokedex_api.Cache, s *GameState) error {

	// c := params.config
	// cache := params.cache

	url, err := c.Previous()

	if err != nil {
		return err
	}

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
	areas, err := pokedex_api.GetAreas(body)

	if err != nil {
		return err
	}

	printAreas(areas)

	return nil
}
