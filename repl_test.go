package main

import "testing"
import "reflect"

func TestCleanInput(t *testing.T) {

	cases := []struct {
		input    string
		expected []string
	}{
		{"  hello  world  ",
			[]string{"hello", "world"}},
	}

	for _, c := range cases {
		result := cleanInput(c.input)
		if !reflect.DeepEqual(result, c.expected) {
			t.Errorf("Expected %v, got %v", c.expected, result)
		}
	}
}
