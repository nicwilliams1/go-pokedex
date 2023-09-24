package main

import (
	"encoding/json"
	"io"
	"os"

	pokedex_api "github.com/nicwilliams1/pokedexcli/internal"
)

type User struct {
	Name    string                         `json:"name"`
	Pokedex map[string]pokedex_api.Pokemon `json:"pokedex"`
}

func (u *User) AddPokemonToUserPokedex(p pokedex_api.Pokemon) {
	_, ok := u.Pokedex[p.Name]
	if !ok {
		u.Pokedex[p.Name] = p
	}
}

type Users struct {
	Users []User `json:"users"`
}

func NewUser(name string, users *Users) User {
	u := User{
		Name:    name,
		Pokedex: make(map[string]pokedex_api.Pokemon),
	}
	users.Users = append(users.Users, u)
	return u
}

func GetUser(name string, users Users) (User, bool) {
	userList := users.Users
	for _, user := range userList {
		if user.Name == name {
			return user, true
		}
	}
	return User{}, false
}

func SaveUsers(users Users) error {
	file, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile("users.json", file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadUsersFromFile() (Users, error) {
	if _, err := os.Stat("users.json"); err == nil {
		var users Users
		file, err := os.Open("users.json")
		if err != nil {
			return users, err
		}
		defer file.Close()

		byteValue, err := io.ReadAll(file)
		if err != nil {
			return users, err
		}

		json.Unmarshal(byteValue, &users)

		return users, nil

	} else {
		// file doesn't exist, return an empty Users struct
		var users Users
		return users, nil
	}
}
