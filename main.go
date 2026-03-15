package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config) error
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
}

func main() {
	cfg := config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		input := scanner.Scan()
		if !input {
			break
		}
		text := cleanInput(scanner.Text())
		if typeofCommand, exists := cliCommands[text[0]]; exists {
			err := typeofCommand.callback(&cfg)
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", text[0])
		}
	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	println("Welcome to the Pokedex!\n")
	println("Usage:")
	println("help: Displays a help message")
	println("map: List the first 20 locations in the Pokemon world")
	println("mapb: List the previous 20 locations in the Pokemon world")
	println("exit: Exit the Pokedex")
	return nil
}

func processLocations(res *http.Response, cfg *config) {
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	var locationAreaResponse locationAreaResponse
	if err := json.Unmarshal(body, &locationAreaResponse); err != nil {
		log.Fatal(err)
	}
	for _, area := range locationAreaResponse.Results {
		fmt.Printf("%s\n", area.Name)
	}
	cfg.Next = locationAreaResponse.Next
	cfg.Previous = locationAreaResponse.Previous
}

func commandMap(cfg *config) error {
	res, err := http.Get(cfg.Next)
	if err != nil {
		log.Fatal(err)
	}
	processLocations(res, cfg)
	return nil
}

func commandMapBack(cfg *config) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	res, err := http.Get(cfg.Previous)
	if err != nil {
		log.Fatal(err)
	}
	processLocations(res, cfg)
	return nil
}
