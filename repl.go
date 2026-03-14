package main

import "strings"

func cleanInput(text string) []string {
	fields := strings.Fields(strings.ToLower(text))
	return fields
}
