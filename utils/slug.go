package utils

import (
	"regexp"
	"strings"
)

func Slug(value []string) string {
	var toSlug string = strings.Join(value, " ")
	toSlug = strings.ToLower(toSlug)
	var re = regexp.MustCompile(`[*+~.()_'":@/]`)
	toSlug = re.ReplaceAllString(toSlug, "")
	toSlug = strings.Replace(toSlug, " ", "-", -1)
	toSlug = strings.Replace(toSlug, "--", "-", -1)
	toSlug = strings.Replace(toSlug, "&", "and", -1)
	return toSlug
}
