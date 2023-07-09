package utils

import (
	"regexp"
	"strings"

	"github.com/gosimple/slug"
)

func Slug(value []string) string {
	var toSlug string = strings.Join(value, " ")
	slugified := slug.Make(toSlug)
	regex := regexp.MustCompile(`[*+~.()_'":@/]`)
	return regex.ReplaceAllString(slugified, "")
}

