package main

import (
    "reflect"
    "testing"
)

func TestCleanInput(t *testing.T) {
    tests := []struct {
        input    string
        expected []string
    }{
        {"hello world", []string{"hello", "world"}},
        {"Charmander Bulbasaur PIKACHU", []string{"charmander", "bulbasaur", "pikachu"}},
        {"   leading and trailing   ", []string{"leading", "and", "trailing"}},
        {"", []string{}},
        {"   ", []string{}},
    }

    for _, test := range tests {
        result := cleanInput(test.input)
        if !reflect.DeepEqual(result, test.expected) {
            t.Errorf("cleanInput(%q) = %v; want %v", test.input, result, test.expected)
        }
    }
}