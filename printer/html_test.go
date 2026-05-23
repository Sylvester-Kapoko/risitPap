package printer

import (
	"strings"
	"testing"
	"time"

	"github.com/Sylvester-Kapoko/risitPap/domain"
	"github.com/shopspring/decimal"
)

func TestHtmlFormatter_CustomerFields(t *testing.T) {
	// A receipt with customer details
	r := &domain.Receipt{
		TransactionID: "INV-001",
		StoreName:     "Test Store",
		StoreAddr:     "123 Main St",
		CustomerName:  "Jane Doe",
		CustomerPhone: "0712345678",
		Currency:      "Ksh",
		Items: []domain.ReceiptItem{
			{Name: "Widget", Qty: 2, UnitPrice: decimal.NewFromFloat(10.00)},
		},
		TaxRate:   decimal.Zero,
		Payment:   domain.Payment{Method: "Cash", Amount: decimal.NewFromInt(20)},
		CreatedAt: time.Now(),
	}

	cfg := HtmlConfig{
		DeveloperName:  "Dev",
		DeveloperPhone: "123",
		ShowFooter:     true,
	}
	fm := NewHtmlFormatter(cfg)
	html, err := fm.Format(r)
	if err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	// Check that customer name and phone appear
	if !strings.Contains(html, "Jane Doe") {
		t.Error("expected customer name 'Jane Doe' in output")
	}
	if !strings.Contains(html, "0712345678") {
		t.Error("expected customer phone '0712345678' in output")
	}

	// Now a receipt WITHOUT customer info
	r2 := &domain.Receipt{
		TransactionID: "INV-002",
		StoreName:     "Test Store",
		StoreAddr:     "123 Main St",
		Currency:      "Ksh",
		Items: []domain.ReceiptItem{
			{Name: "Widget", Qty: 1, UnitPrice: decimal.NewFromFloat(5.00)},
		},
		TaxRate:   decimal.Zero,
		Payment:   domain.Payment{Method: "Cash", Amount: decimal.NewFromInt(10)},
		CreatedAt: time.Now(),
	}

	html2, err := fm.Format(r2)
	if err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	// The phrase "To:" should NOT appear when no customer
	if strings.Contains(html2, "To:") {
		t.Error("did NOT expect 'To:' in output when customer is absent")
	}
}
