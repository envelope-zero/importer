package util

import (
	"strings"

	"github.com/shopspring/decimal"
)

// NormalizeGerman takes in a German formatted numer
// and returns a decimal
func NormalizeGerman(i string) (decimal.Decimal, error) {
	s := strings.Replace(i, ".", "", -1)
	s = strings.Replace(s, ",", ".", 1)
	return decimal.NewFromString(s)
}
