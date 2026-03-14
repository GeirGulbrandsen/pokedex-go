package main

import "testing"
import "reflect"

func TestCleanInput(t *testing.T) {
	input := "  Hello, World!  "
	expected := []string{}
	result := cleanInput(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
