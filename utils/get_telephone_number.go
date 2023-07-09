package utils

import "regexp"

func GetTelephoneNumber(s string) string {
	re := regexp.MustCompile("[^0-9]+")
	return re.ReplaceAllString(s, "")
}

func TransformTelephoneNumber(s string) string {
	total := len(s)

	modal := total % 10

	switch modal {
	case 0:
		return s[0:10]
	case 1:
		if total == 11 {
			return s[0:11]
		}
		return s[0:10]
	case 2:
		if total == 22 {
			return s[0:11]
		}
		return s[0:10]
	case 3:
		if total == 33 {
			return s[0:11]
		}
		return s[0:10]
	default:
		return s[0:10]
	}
}
