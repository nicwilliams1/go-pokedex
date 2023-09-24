package main

import (
	"fmt"
	"math/rand"
	"slices"

	pokedex_api "github.com/nicwilliams1/pokedexcli/internal"
)

func commandCatch(c *config, cache *pokedex_api.Cache, s *GameState) error {
	if !slices.Contains(s.area.pokemonNames, s.catchAttempt) {
		return fmt.Errorf("pokemon [%s] not present in current area", s.catchAttempt)
	}

	if _, ok := s.user.Pokedex[s.catchAttempt]; ok {
		return fmt.Errorf("you've already caught pokemon [%s]", s.catchAttempt)
	}

	fmt.Printf("Throwing a Pokeball at %s... \n", s.catchAttempt)

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", s.catchAttempt)
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
	pokemon, err := pokedex_api.GetPokemon(body)

	if err != nil {
		return err
	}

	// get experience level of pokemon
	base_experience := pokemon.BaseExperience
	roll := rand.Intn(100)
	minRoll := int(float64(base_experience) * 0.75)

	if roll >= 95 || roll > minRoll {
		fmt.Printf("[roll %v] %s was caught!\n", roll, pokemon.Name)
		fmt.Println("You may now inspect it using the inspect command")

		s.user.AddPokemonToUserPokedex(pokemon)

		return nil
	}

	fmt.Printf("[roll %v, %v needed] %s escaped!\n", roll, minRoll, pokemon.Name)

	return nil
}
