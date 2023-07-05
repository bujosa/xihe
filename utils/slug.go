package utils

import "strings"

func Slug(value string) string {
	return strings.ReplaceAll(strings.ToLower(value), " ", "-")
}

