package pokedex_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type PokedexArea struct {
	Name string
	Url  string
}

type ApiResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type ApiResponseArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetApiResponse(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed", url)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode == 404 {
		return nil, errors.New("resource not found")
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("request to %s returned with status code %d", url, res.StatusCode)
	}
	if err != nil {
		return nil, fmt.Errorf("request to %s failed", url)
	}

	return body, nil
}

func GetAreas(body []byte) ([]PokedexArea, error) {
	var result ApiResponse
	err := json.Unmarshal(body, &result)

	if err != nil {
		return nil, errors.New("could not unmarshal JSON")
	}

	areas := make([]PokedexArea, 0)

	for _, area := range result.Results {
		pArea := PokedexArea{Name: area.Name, Url: area.URL}
		areas = append(areas, pArea)
	}

	return areas, nil
}

func GetPokemonInArea(body []byte) ([]string, error) {
	var result ApiResponseArea
	err := json.Unmarshal(body, &result)

	if err != nil {
		return nil, errors.New("could not unmarshal JSON")
	}

	pokemon := make([]string, 0)

	for _, row := range result.PokemonEncounters {
		pokemon = append(pokemon, row.Pokemon.Name)
	}

	return pokemon, nil

}

func GetPokemon(body []byte) (Pokemon, error) {
	var result Pokemon
	err := json.Unmarshal(body, &result)

	if err != nil {
		var p Pokemon
		return p, errors.New("could not unmarshal JSON")
	}

	return result, nil

}
