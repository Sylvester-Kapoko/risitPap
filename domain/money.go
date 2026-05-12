package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Money wraps a decimal for currency display.
type Money struct {
	Amount decimal.Decimal
}


func DisplayCurrency(amount decimal.Decimal, currency string) string {
    if currency == "" {
        currency = "Ksh"   // default for your market
    }
    return fmt.Sprintf("%s%s", currency, amount.StringFixed(2))
}



