package domain

import (
	"time"
	"github.com/shopspring/decimal"

)

type ReceiptItem struct {
	Name      string
	Qty       int
	UnitPrice decimal.Decimal
}

type Receipt struct {
	TransactionID        string
	StoreName string
	StoreAddr string
	StorePhone string
	StoreTaxID string
	HasVAT bool
	Currency string 
	Items     []ReceiptItem
	TaxRate   decimal.Decimal
	Payment Payment
	CreatedAt time.Time
}


func (r *Receipt) Subtotal() decimal.Decimal {
    sum := decimal.Zero
    for _, item := range r.Items {
        sum = sum.Add(item.LineTotal())
    }
    return sum
}

func (r *Receipt) Tax() decimal.Decimal {
    return r.Subtotal().Mul(r.TaxRate)
}

func (r *Receipt) Total() decimal.Decimal {
    return r.Subtotal().Add(r.Tax())
}

func (r *Receipt) Change() decimal.Decimal {
    return r.Payment.Amount.Sub(r.Total())
}

func (r *Receipt) TaxRatePct() string{
     return r.TaxRate.Mul(decimal.NewFromInt(100)).StringFixed(2)
}

// LineTotal returns the total for this line ( qty x unit price).

func (item *ReceiptItem) LineTotal() decimal.Decimal {
    return item.UnitPrice.Mul(decimal.NewFromInt(int64(item.Qty)))
}

// Payment records how the customer paid
type Payment struct {
     Method string
     Amount decimal.Decimal

}

