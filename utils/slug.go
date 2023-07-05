package utils

import "strings"

func Slug(value []string) string {
	// Unified string with - as separator
	var unified string = strings.Join(value, "-")
	return strings.ReplaceAll(strings.ToLower(unified), " ", "-")
}

