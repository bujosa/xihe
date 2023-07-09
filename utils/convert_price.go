package utils

func ConvertPrice(price int, currency string) int {
	if currency == "dop" {
		return price / DOP_TO_USD
	}
	return price
}
