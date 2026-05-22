package main

import (
	"os"
	"time"

	"github.com/Sylvester-Kapoko/risitPap/domain"
	"github.com/Sylvester-Kapoko/risitPap/printer"
	"github.com/shopspring/decimal"
)

func main() {
	receipt := &domain.Receipt{
		StoreName: "BOB'S CAFE",
		StoreAddr: "123 Main St",
		Items: []domain.ReceiptItem{
			{Name: "Burger", Qty: 2, UnitPrice: decimal.RequireFromString("12.00")},
			{Name: "Fries", Qty: 1, UnitPrice: decimal.RequireFromString("5.00")},
		},
		TaxRate: decimal.RequireFromString("0.10"),
		Payment: domain.Payment{
			Method: "Cash",
			Amount: decimal.RequireFromString("40.00"),
		},
		CreatedAt: time.Now(),
	}

	formatter := printer.NewPlainTextFormatter(32)
	receiptPrinter := printer.NewReceiptPrinter(formatter, os.Stdout)

	if err := receiptPrinter.Print(receipt); err != nil {
		panic(err)
	}
}
