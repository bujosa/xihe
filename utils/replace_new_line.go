package utils

import "strings"

func ReplaceNewLine(s string) string {
	return strings.Replace(s, "\n", " ", -1)
}
