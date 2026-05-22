package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Money wraps a decimal for currency display.
type Money struct {
	Amount decimal.Decimal
}

// DisplayCurrency formats a decimal amount with a currency symbol.
//
// Precondition: amount is any valid decimal.Decimal (including zero).
// Postcondition: the result has exactly two decimal places.
// Invariant: for any two currencies c1, c2,
//   numericPart(DisplayCurrency(amount, c1)) == numericPart(DisplayCurrency(amount, c2))
//
// If currency is "", "Ksh" is used as the default.

func DisplayCurrency(amount decimal.Decimal, currency string) string {
	if currency == "" {
		currency = "Ksh" // default for your market
	}
	return fmt.Sprintf("%s%s", currency, amount.StringFixed(2))
}
