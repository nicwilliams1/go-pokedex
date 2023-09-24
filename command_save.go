package main

import (
	"errors"

	pokedex_api "github.com/nicwilliams1/pokedexcli/internal"
)

func commandSave(c *config, cache *pokedex_api.Cache, s *GameState) error {

	if s.user.Name == "guest" {
		return errors.New("guest user cannot save state")
	}

	err := SaveUsers(s.users)
	if err != nil {
		return errors.New("failed to save to file")
	}
	return nil
}
