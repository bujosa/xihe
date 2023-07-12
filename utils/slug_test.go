package utils

import (
	"log"
	"testing"
)

func TestSlug(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"Hello World", "hello-world"},
		{"This is a Test", "this-is-a-test"},
		{"Guelo Motors, SRL", "guelo-motors,-srl"},
		{"", ""},
		{"Ismael Cruz Automóviles", "ismael-cruz-automoviles"},
		{"Gold`S Brothers Auto ", "gold`s-brothers-auto"},
		{"Piña´s Auto", "pina´s-auto"},
	}

	for _, test := range tests {
		log.Print(test.input + "\n")

		convert := []string{test.input}
		result := Slug(convert)
		if result != test.output {
			t.Errorf("Slug(%q) = %q, expected %q", test.input, result, test.output)
		}
	}
}
