package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/geirgulbrandsen/pokedex-go/internal/pokecache"
)

type config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
	Pokedex  map[string]pokemonData
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args []string) error
}

type locationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type locationAreaResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []locationArea `json:"results"`
}

type pokemon struct {
	Name string `json:"name"`
}

type pokemonData struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

type pokemonEncounter struct {
	Pokemon pokemon `json:"pokemon"`
}

type locationAreaDetailsResponse struct {
	PokemonEncounters []pokemonEncounter `json:"pokemon_encounters"`
}

var cliCommands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Show available commands",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Show the first 20 locations in the Pokemon world",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Show the previous 20 locations in the Pokemon world",
		callback:    commandMapBack,
	},
	"explore": {
		name:        "explore",
		description: "Explore a location area for pokemon encounters",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "Attempt to catch a Pokemon by name",
		callback:    commandCatch,
	},
}

func main() {
	cfg := config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		Cache:    pokecache.NewCache(5 * time.Second),
		Pokedex:  make(map[string]pokemonData),
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		input := scanner.Scan()
		if !input {
			break
		}
		text := cleanInput(scanner.Text())
		if len(text) == 0 {
			continue
		}
		if typeofCommand, exists := cliCommands[text[0]]; exists {
			err := typeofCommand.callback(&cfg, text[1:])
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", text[0])
		}
	}
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	println("Welcome to the Pokedex!\n")
	println("Usage:")
	println("help: Displays a help message")
	println("map: List the first 20 locations in the Pokemon world")
	println("mapb: List the previous 20 locations in the Pokemon world")
	println("explore <area_name>: Lists pokemon in a location area")
	println("catch <pokemon_name>: Attempt to catch a Pokemon")
	println("exit: Exit the Pokedex")
	return nil
}

func processLocations(body []byte, cfg *config) error {
	var locationAreaResponse locationAreaResponse
	if err := json.Unmarshal(body, &locationAreaResponse); err != nil {
		return err
	}
	for _, area := range locationAreaResponse.Results {
		fmt.Printf("%s\n", area.Name)
	}
	cfg.Next = locationAreaResponse.Next
	cfg.Previous = locationAreaResponse.Previous
	return nil
}

func getURLData(cfg *config, url string) ([]byte, error) {
	if cachedData, ok := cfg.Cache.Get(url); ok {
		return cachedData, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
	}

	cfg.Cache.Add(url, body)
	return body, nil
}

func commandMap(cfg *config, args []string) error {
	body, err := getURLData(cfg, cfg.Next)
	if err != nil {
		return err
	}
	return processLocations(body, cfg)
}

func commandMapBack(cfg *config, args []string) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	body, err := getURLData(cfg, cfg.Previous)
	if err != nil {
		return err
	}
	return processLocations(body, cfg)
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a location area name")
	}

	areaName := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", areaName)

	body, err := getURLData(cfg, url)
	if err != nil {
		return err
	}

	var details locationAreaDetailsResponse
	if err := json.Unmarshal(body, &details); err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", areaName)
	fmt.Println("Found Pokemon:")
	for _, encounter := range details.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a pokemon name")
	}

	name := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)

	body, err := getURLData(cfg, url)
	if err != nil {
		return err
	}

	var p pokemonData
	if err := json.Unmarshal(body, &p); err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	baseExp := p.BaseExperience
	if baseExp <= 0 {
		baseExp = 1
	}

	if rand.Intn(baseExp) < 40 {
		fmt.Printf("%s was caught!\n", name)
		cfg.Pokedex[name] = p
	} else {
		fmt.Printf("%s escaped!\n", name)
	}

	return nil
}
