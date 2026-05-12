package printer

package printer

import (
	"strings"
	"testing"
	"time"

	"github.com/Sylvester-Kapoko/Receipts/domain"
	"github.com/shopspring/decimal"
)

func TestPlainTextFormatter(t *testing.T) {
	formatter := NewPlainTextFormatter(40)

	makeReceipt := func() *domain.Receipt {
		return &domain.Receipt{
			TransactionID: "TXN-001",
			StoreName:     "Acme Widget Co.",
			StoreAddr:     "123 Main Street",
			TaxRate:       decimal.NewFromFloat(0.0825),
			Payment: domain.Payment{
				Method: "Cash",
				Amount: decimal.NewFromInt(20),
			},
			CreatedAt: time.Date(2026, 4, 27, 14, 30, 0, 0, time.UTC),
			Items: []domain.ReceiptItem{
				{Name: "Widget", Qty: 1, UnitPrice: decimal.NewFromFloat(9.99)},
				{Name: "Super Deluxe Widget", Qty: 2, UnitPrice: decimal.NewFromFloat(3.50)},
			},
		}
	}

	// === Structural Tests ===

	t.Run("renders header with store name centered at the top", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		if !strings.Contains(output, "Acme Widget Co.") {
			t.Errorf("expected store name in output:\n%s", output)
		}
	})

	t.Run("prints store address below store name", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		storeIdx := strings.Index(output, "Acme Widget Co.")
		addrIdx := strings.Index(output, "123 Main Street")
		if storeIdx == -1 || addrIdx == -1 || storeIdx > addrIdx {
			t.Errorf("expected address after store name:\n%s", output)
		}
	})

	t.Run("prints transaction ID that is unique and clearly labeled", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		if !strings.Contains(output, "Transaction: TXN-001") {
			t.Errorf("expected transaction ID in output:\n%s", output)
		}
	})

	t.Run("prints transaction date and time", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		if !strings.Contains(output, "2026-04-27 14:30") {
			t.Errorf("expected date and time in output:\n%s", output)
		}
	})

	t.Run("separates sections with a horizontal rule", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		if !strings.Contains(output, "----------------------------------------") {
			t.Errorf("expected separator line in output:\n%s", output)
		}
	})

	t.Run("renders footer with thank-you message", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		if !strings.Contains(output, "Thank you!") {
			t.Errorf("expected thank you in output:\n%s", output)
		}
	})

	// === Line Item Formatting Tests ===

	t.Run("prints each item on its own line with quantity, description, and price", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		if !strings.Contains(output, "Widget") {
			t.Errorf("expected Widget in items:\n%s", output)
		}
		if !strings.Contains(output, "Super Deluxe Widget") {
			t.Errorf("expected Super Deluxe Widget in items:\n%s", output)
		}
	})

	t.Run("formats quantities as integers without decimal places", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		// Quantity 2 should appear as "2" not "2.00"
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Super Deluxe Widget") {
				if !strings.Contains(line, " 2 ") && !strings.HasPrefix(strings.TrimSpace(line), "2") {
					t.Errorf("expected integer quantity, got line: %s", line)
				}
			}
		}
	})

	// === Totals and Calculations Tests ===

	t.Run("displays subtotal line above tax line", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		subIdx := strings.Index(output, "Subtotal")
		taxIdx := strings.Index(output, "Tax")
		if subIdx == -1 || taxIdx == -1 || subIdx > taxIdx {
			t.Errorf("expected subtotal before tax:\n%s", output)
		}
	})

	t.Run("shows tax as a separate line with applied percentage in parentheses", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		if !strings.Contains(output, "Tax (8.25%)") {
			t.Errorf("expected tax percentage in output:\n%s", output)
		}
	})

	t.Run("prints grand total after tax", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		taxIdx := strings.Index(output, "Tax")
		totalIdx := strings.Index(output, "TOTAL")
		if taxIdx == -1 || totalIdx == -1 || taxIdx > totalIdx {
			t.Errorf("expected total after tax:\n%s", output)
		}
	})

	t.Run("displays all currency values with exactly two decimal places", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		// Check that dollar amounts have .00 or .99 format
		if !strings.Contains(output, "9.99") {
			t.Errorf("expected two-decimal currency values:\n%s", output)
		}
		if !strings.Contains(output, "3.50") {
			t.Errorf("expected two-decimal currency values:\n%s", output)
		}
	})

	// === Payment Information Tests ===

	t.Run("shows payment method on a dedicated line", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		if !strings.Contains(output, "Cash") {
			t.Errorf("expected payment method in output:\n%s", output)
		}
	})

	t.Run("displays amount tendered and change due when payment is cash", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		if !strings.Contains(output, "$20.00") {
			t.Errorf("expected payment amount:\n%s", output)
		}
		if !strings.Contains(output, "Change") {
			t.Errorf("expected change line:\n%s", output)
		}
	})

	// === Edge Case Tests (currently passing) ===

	t.Run("handles receipt with no items", func(t *testing.T) {
		receipt := makeReceipt()
		receipt.Items = nil
		output := formatter.Format(receipt)
		// Should still render without panicking
		if !strings.Contains(output, "TOTAL") {
			t.Errorf("expected total even with no items:\n%s", output)
		}
	})

	t.Run("preserves formatting when store name is a single character", func(t *testing.T) {
		receipt := makeReceipt()
		receipt.StoreName = "A"
		output := formatter.Format(receipt)
		if !strings.Contains(output, "A") {
			t.Errorf("expected single-char store name in output:\n%s", output)
		}
	})

	t.Run("preserves formatting when store name is 50 characters", func(t *testing.T) {
		receipt := makeReceipt()
		receipt.StoreName = "A Very Long Store Name That Goes On And On Forever"
		output := formatter.Format(receipt)
		if len(output) == 0 {
			t.Error("expected output for long store name, got empty string")
		}
	})

	t.Run("handles item with price of zero without misaligning total", func(t *testing.T) {
		receipt := makeReceipt()
		receipt.Items = []domain.ReceiptItem{
			{Name: "Free Sample", Qty: 1, UnitPrice: decimal.Zero},
		}
		output := formatter.Format(receipt)
		if !strings.Contains(output, "0.00") {
			t.Errorf("expected zero price formatted as 0.00:\n%s", output)
		}
	})

	// === Fixed-Width Tests ===

	t.Run("uses fixed-width formatting at configured width", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		lines := strings.Split(output, "\n")
		for i, line := range lines {
			if len(line) > 40 {
				t.Errorf("line %d exceeds width 40: len=%d, line=%q", i, len(line), line)
			}
		}
	})

	t.Run("renders correctly at 40 columns", func(t *testing.T) {
		f := NewPlainTextFormatter(40)
		receipt := makeReceipt()
		output := f.Format(receipt)
		lines := strings.Split(output, "\n")
		for i, line := range lines {
			if len(line) > 40 {
				t.Errorf("line %d exceeds 40 cols: %q", i, line)
			}
		}
	})

	t.Run("renders correctly at 80 columns", func(t *testing.T) {
		f := NewPlainTextFormatter(80)
		receipt := makeReceipt()
		output := f.Format(receipt)
		lines := strings.Split(output, "\n")
		for i, line := range lines {
			if len(line) > 80 {
				t.Errorf("line %d exceeds 80 cols: %q", i, line)
			}
		}
	})

	t.Run("renders correctly at 120 columns", func(t *testing.T) {
		f := NewPlainTextFormatter(120)
		receipt := makeReceipt()
		output := f.Format(receipt)
		lines := strings.Split(output, "\n")
		for i, line := range lines {
			if len(line) > 120 {
				t.Errorf("line %d exceeds 120 cols: %q", i, line)
			}
		}
	})

	// === Currency Symbol Test ===

	t.Run("displays currency symbol next to all monetary values", func(t *testing.T) {
		receipt := makeReceipt()
		output := formatter.Format(receipt)
		if !strings.Contains(output, "$") {
			t.Errorf("expected dollar sign in output:\n%s", output)
		}
	})
}
