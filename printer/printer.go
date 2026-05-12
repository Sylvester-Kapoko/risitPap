package printer

import (
	"io"

	"github.com/Sylvester-Kapoko/Receipts/domain"
)

// ReceiptPrinter sends a formatted receipt to a writer.
type ReceiptPrinter struct {
	formatter ReceiptFormatter
	writer    io.Writer
}

// NewReceiptPrinter creates a printer with the given formatter and output destination.
func NewReceiptPrinter(f ReceiptFormatter, w io.Writer) *ReceiptPrinter {
	return &ReceiptPrinter{formatter: f, writer: w}
}

// Print formats the receipt and writes it to the output.
func (p *ReceiptPrinter) Print(receipt *domain.Receipt) error {
	output := p.formatter.Format(receipt)
	_, err := io.WriteString(p.writer, output)
	return err
}
