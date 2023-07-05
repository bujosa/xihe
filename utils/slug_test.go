package utils

import "testing"

func TestSlug(t *testing.T) {
    tests := []struct {
        input  string
        output string
    }{
        {"Hello World", "hello-world"},
        {"This is a Test", "this-is-a-test"},
        {"", ""},
    }

    for _, test := range tests {
        result := Slug(test.input)
        if result != test.output {
            t.Errorf("Slug(%q) = %q, expected %q", test.input, result, test.output)
        }
    }
}