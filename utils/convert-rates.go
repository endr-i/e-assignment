package utils

import "github.com/shopspring/decimal"

func ConvertRate(rate1 decimal.Decimal, rate2 decimal.Decimal) decimal.Decimal {
	return rate1.Div(rate2) // conversion algo might be more complicated
}
