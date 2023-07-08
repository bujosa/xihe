package utils

import "regexp"

func GetTelephoneNumber(s string) string {
    re := regexp.MustCompile("[^0-9]+")
    return re.ReplaceAllString(s, "")
}