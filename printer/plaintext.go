package printer

import (
	"fmt"
	"strings"

	"github.com/Sylvester-Kapoko/risitPap/domain"
)

type PlainTextFormatter struct {
	width int
}

func NewPlainTextFormatter(width int) *PlainTextFormatter {
	if width < 20 {
		width = 32
	}
	return &PlainTextFormatter{width: width}
}

func (f *PlainTextFormatter) Format(r *domain.Receipt) (string, error) {
	var b strings.Builder

	center(&b, r.StoreName, f.width)
	center(&b, r.StoreAddr, f.width)
	center(&b, fmt.Sprintf("Transaction: %s", r.TransactionID), f.width)
	b.WriteString("\n")

	for _, item := range r.Items {
		itemLine(&b, item.Name, fmt.Sprintf("%d", item.Qty),
			domain.DisplayCurrency(item.UnitPrice, r.Currency), domain.DisplayCurrency(item.LineTotal(), r.Currency), f.width)
	}

	separator(&b, f.width)
	moneyLine(&b, "Subtotal", domain.DisplayCurrency(r.Subtotal(), r.Currency), f.width)
	if r.HasVAT {
		moneyLine(&b, fmt.Sprintf("Tax (%s%%)", r.TaxRatePct()), domain.DisplayCurrency(r.Tax(), r.Currency), f.width)
	}
	moneyLine(&b, "TOTAL", domain.DisplayCurrency(r.Total(), r.Currency), f.width)
	b.WriteString("\n")
	moneyLine(&b, r.Payment.Method, domain.DisplayCurrency(r.Payment.Amount, r.Currency), f.width)
	moneyLine(&b, "Change", domain.DisplayCurrency(r.Change(), r.Currency), f.width)

	b.WriteString("\n")
	center(&b, "Thank you!", f.width)
	center(&b, r.CreatedAt.Format("2006-01-02 15:04"), f.width)

	return b.String(), nil
}
