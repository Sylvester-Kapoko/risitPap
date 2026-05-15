// domain/receipt.go
package domain

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type ReceiptItem struct {
	Name      string
	Qty       int
	UnitPrice decimal.Decimal
}

type Receipt struct {
	ID            string        `json:"id"`
	TransactionID string        `json:"transaction_id"`
	StoreName     string
	StoreAddr     string
	StorePhone    string
	StoreTaxID    string
	HasVAT        bool
	Currency      string
	Items         []ReceiptItem
	TaxRate       decimal.Decimal
	Payment       Payment
	Status        PaymentStatus `json:"status"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

func NewReceipt() *Receipt {
	now := time.Now().UTC()
	return &Receipt{
		Status:    StatusPending,
		TaxRate:   decimal.Zero,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ==================== PAYMENT STATUS (v1 Feature) ====================

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "PENDING"
	StatusPaid      PaymentStatus = "PAID"
	StatusFailed    PaymentStatus = "FAILED"
	StatusCancelled PaymentStatus = "CANCELLED"
	StatusRefunded  PaymentStatus = "REFUNDED"
)

func ValidPaymentStatusTransition(from, to PaymentStatus) bool {
	switch from {
	case StatusPending:
		return to == StatusPaid || to == StatusFailed || to == StatusCancelled
	case StatusPaid:
		return to == StatusRefunded
	default:
		return false
	}
}

func ValidateStatusTransition(from, to PaymentStatus) error {
	if from == to {
		return nil
	}
	if ValidPaymentStatusTransition(from, to) {
		return nil
	}
	return errors.New("invalid payment status transition: " + string(from) + " → " + string(to))
}

func (r *Receipt) ChangeStatus(newStatus PaymentStatus, reason string) error {
	if err := ValidateStatusTransition(r.Status, newStatus); err != nil {
		return err
	}
	r.Status = newStatus
	r.UpdatedAt = time.Now().UTC()
	return nil
}

// ==================== PAYMENT STATE ====================

type PaymentState string

const (
	PaymentPending  PaymentState = "PENDING"
	PaymentPartial  PaymentState = "PARTIAL"
	PaymentPaid     PaymentState = "PAID"
	PaymentOverpaid PaymentState = "OVERPAID"
)

// PaymentState returns the current payment status.
// For simplicity, Cash is always Paid; other methods are Pending until confirmed.
func (r *Receipt) PaymentState() PaymentStatus {
    if r.Payment.Status == "" {
        return StatusPending
    }
    return r.Payment.Status
}

func (r *Receipt) AmountState() PaymentState {
	total := r.Total()
	paid := r.Payment.Amount

	if total.IsNegative() {
		return PaymentPending
	}

	switch {
	case paid.IsZero():
		return PaymentPending
	case paid.GreaterThan(total):
		return PaymentOverpaid
	case paid.Equal(total):
		return PaymentPaid
	case paid.LessThan(total) && paid.GreaterThan(decimal.Zero):
		return PaymentPartial
	default:
		return PaymentPending
	}
}

func (r *Receipt) IsFullyPaid() bool {
	state := r.AmountState()
	return state == PaymentPaid || state == PaymentOverpaid
}

func (r *Receipt) IsPending() bool {
	return r.AmountState() == PaymentPending
}

func (r *Receipt) IsPartiallyPaid() bool {
	return r.AmountState() == PaymentPartial
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

func (r *Receipt) TaxRatePct() string {
	return r.TaxRate.Mul(decimal.NewFromInt(100)).StringFixed(2)
}

func (item *ReceiptItem) LineTotal() decimal.Decimal {
	return item.UnitPrice.Mul(decimal.NewFromInt(int64(item.Qty)))
}

type Payment struct {
	Method string
	Amount decimal.Decimal
	Status PaymentStatus
}