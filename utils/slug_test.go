package utils

import "testing"

func TestSlug(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"Hello World", "hello-world"},
		{"This is a Test", "this-is-a-test"},
		{"Guelo Motors, SRL", "guelo-motors,-srl"},
		{"", ""},
	}

	for _, test := range tests {
		print(test.input + "\n")

		convert := []string{test.input}
		result := Slug(convert)
		if result != test.output {
			t.Errorf("Slug(%q) = %q, expected %q", test.input, result, test.output)
		}
	}
}
