package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	pokedex_api "github.com/nicwilliams1/pokedexcli/internal"
)

type GameState struct {
	replState     int
	user          User
	area          Area
	catchAttempt  string
	inspectTarget string
	users         Users
}

func (g *GameState) SetGameState(s int) {
	g.replState = s
}

func (g *GameState) SetCurrentUser(u User) {
	g.user = u
}

func (g *GameState) SetCurrentCatchAttempt(c string) {
	g.catchAttempt = c
}
func (g *GameState) SetCurrentInspectTarget(i string) {
	g.inspectTarget = i
}

type Area struct {
	name         string
	pokemonNames []string
}

func (a *Area) SetCurrentArea(name string, pokemonNames []string) {
	a.name = name
	a.pokemonNames = pokemonNames
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *pokedex_api.Cache, *GameState) error
}

type config struct {
	limit   int
	offset  int
	baseurl string
}

func startRepl() {

	// load users
	users, err := LoadUsersFromFile()
	if err != nil {
		fmt.Println("Failed to load user list")
	}

	// initialize game state obj
	var State = GameState{
		replState: 0,
		user: User{
			Name:    "guest",
			Pokedex: make(map[string]pokedex_api.Pokemon, 0),
		},
		area: Area{
			name:         "",
			pokemonNames: make([]string, 0),
		},
		catchAttempt: "",
		users:        users,
	}

	// load config
	c := config{
		baseurl: "https://pokeapi.co/api/v2/location-area",
		limit:   20,
		offset:  0,
	}

	// load cache
	cache := pokedex_api.NewCache(5 * time.Minute)

	// start input loop
	reader := bufio.NewScanner(os.Stdin)
	for {
		if State.replState == 0 || State.replState == 2 {
			fmt.Printf("Pokedex [%s] >\n", State.user.Name)
		}

		if State.replState == 0 {
			fmt.Println("Welcome guest. Do you wish to sign in? [Y]/[N] >")
		}

		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		if State.replState == 0 {
			// this should be a Y or N to either sign in or continue as a guest
			if words[0] == "Y" || words[0] == "y" {
				fmt.Println("Please enter username to sign in >")
				State.SetGameState(1)
			} else {
				fmt.Println("Continuing as guest")
				State.SetGameState(2)
			}
			continue
		}

		if State.replState == 1 {
			// this is the users name
			username := words[0]
			u, ok := GetUser(username, users)

			if !ok {
				fmt.Printf("Username not found, creating new user: %s\n", username)
				u = NewUser(username, &users)
				err := SaveUsers(users)
				if err != nil {
					fmt.Printf("Failed to create user, continuing as guest")
					State.SetGameState(2)
					continue
				}
			} else {
				fmt.Printf("Welcome %s!\n", u.Name)
			}

			State.user.setUser(u.Name, u.Pokedex)
			State.SetGameState(2)
			continue
		}

		if State.replState == 2 {
			command, ok := getCommands()[words[0]]
			if !ok {
				fmt.Println("Unknown command")
				continue
			}

			if command.name == "explore" {
				if len(words) < 2 {
					fmt.Println("Explore command requires an area name as the second parameter. Type help for more info.")
					continue
				}
				State.area.SetCurrentArea(words[1], make([]string, 0))
			}

			if command.name == "catch" {
				if len(words) < 2 {
					fmt.Println("Catch command requires a pokemon name as the second parameter. Type help for more info.")
					continue
				}
				if State.area.name == "" {
					fmt.Println("You must explore an area using the explore command before you can attempt to catch a pokemon from that area")
					continue
				}

				State.SetCurrentCatchAttempt(words[1])
			}

			if command.name == "inspect" {
				if len(words) < 2 {
					fmt.Println("Inspect command requires a pokemon name as the second parameter. Type help for more info.")
					continue
				}
				State.SetCurrentInspectTarget(words[1])
			}

			err := command.callback(&c, cache, &State)

			if err != nil {
				fmt.Println(err)
			}
		}

	}
}

func (u *User) setUser(name string, pokedex map[string]pokedex_api.Pokemon) {
	u.Name = name
	u.Pokedex = pokedex
}

func cleanInput(text string) []string {
	lowerCase := strings.ToLower(text)
	words := strings.Fields(lowerCase)
	return words
}

func (c *config) Next() string {
	c.offset += c.limit
	return fmt.Sprintf("%s?offset=%v&limit=%v", c.baseurl, c.offset-20, c.limit)
}
func (c *config) Previous() (string, error) {
	t := c.offset
	if t-c.limit <= 0 {
		return "", errors.New("no previous map pages available")
	}

	c.offset -= c.limit
	return fmt.Sprintf("%s?offset=%v&limit=%v", c.baseurl, c.offset-20, c.limit), nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"save": {
			name:        "save",
			description: "Saves current progress to file",
			callback:    commandSave,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 pokemon world locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 pokemon world locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays pokemon found in given area. Requires [area name] as parameter",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a pokemon in your current area. Requires [pokemon name] as a parameter",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Display info about a pokemon you have already caught. Required [pokemon name] as a parameter",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays pokemon you have caught",
			callback:    commandPokedex,
		},
	}
}
