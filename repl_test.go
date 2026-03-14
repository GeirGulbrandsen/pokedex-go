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
		{"   foo   bar   baz   ",
			[]string{"foo", "bar", "baz"}},
		{"   singleword   ",
			[]string{"singleword"}},
		{"   multiple    spaces   between   words   ",
			[]string{"multiple", "spaces", "between", "words"}},
	}

	for _, c := range cases {
		result := cleanInput(c.input)
		if !reflect.DeepEqual(result, c.expected) {
			t.Errorf("Expected %v, got %v", c.expected, result)
		}
	}
}
