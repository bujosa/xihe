package utils

import (
	"regexp"
	"strings"
)

func Slug(value []string) string {
	var toSlug string = strings.Join(value, " ")
	toSlug = strings.ToLower(toSlug)
	toSlug = strings.Replace(toSlug, "ñ", "n", -1)
	toSlug = ReplaceAccentedChars(toSlug)
	toSlug = strings.TrimSpace(toSlug)
	var re = regexp.MustCompile(`[*+~.()_'":@/]`)
	toSlug = re.ReplaceAllString(toSlug, "")
	toSlug = strings.Replace(toSlug, " ", "-", -1)
	toSlug = strings.Replace(toSlug, "--", "-", -1)
	toSlug = strings.Replace(toSlug, "&", "and", -1)
	return toSlug
}

func ReplaceAccentedChars(value string) string {
	value = strings.Replace(value, "á", "a", -1)
	value = strings.Replace(value, "é", "e", -1)
	value = strings.Replace(value, "í", "i", -1)
	value = strings.Replace(value, "ó", "o", -1)
	value = strings.Replace(value, "ú", "u", -1)
	value = strings.Replace(value, "Á", "A", -1)
	value = strings.Replace(value, "É", "E", -1)
	value = strings.Replace(value, "Í", "I", -1)
	value = strings.Replace(value, "Ó", "O", -1)
	value = strings.Replace(value, "Ú", "U", -1)
	return value
}
