package domain

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestNewReceipt(t *testing.T) {
	r := NewReceipt()
	if r.Status != StatusPending {
		t.Errorf("expected PENDING, got %s", r.Status)
	}
}

func TestPaymentStateCalculation(t *testing.T) {
	tests := []struct {
		name     string
		total    string
		paid     string
		expected PaymentState
	}{
		{"nothing paid", "1000", "0", PaymentPending},
		{"partial", "1000", "300", PaymentPartial},
		{"exact", "1000", "1000", PaymentPaid},
		{"overpaid", "1000", "1200", PaymentOverpaid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReceipt()
			r.Items = []ReceiptItem{
				{Name: "Test", Qty: 1, UnitPrice: decimal.RequireFromString(tt.total)},
			}
			r.TaxRate = decimal.Zero
			r.Payment.Amount = decimal.RequireFromString(tt.paid)

			if got := r.AmountState(); got != tt.expected {
				t.Errorf("got %s, want %s", got, tt.expected)
			}
		})
	}
}
