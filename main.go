package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		input := scanner.Scan()
		if !input {
			break
		}
		text := cleanInput(scanner.Text())
		if typeofCommand, exists := cliCommands[text[0]]; exists {
			err := typeofCommand.callback()
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", text[0])
		}
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	println("Welcome to the Pokedex!\n")
	println("Usage:")
	println("help: Displays a help message")
	println("exit: Exit the Pokedex")
	return nil
}
