package utils

import (
	"strings"
)

func Title(s string) string {
	words := strings.Split(s, " ")
	for i, word := range words {
		// First letter to upper case
		words[i] = capitalizeFirstLetter(word)
	}
	return strings.Join(words, " ")
}

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
