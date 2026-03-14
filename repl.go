package main

import "strings"

func cleanInput(text string) []string {
	fields := strings.Fields(text)
	return fields
}
