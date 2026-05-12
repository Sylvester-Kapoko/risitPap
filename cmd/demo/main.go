package main

import (
	
                 "os"
	"time"
	"github.com/shopspring/decimal"
	"github.com/Sylvester-Kapoko/Receipts/domain"
	"github.com/Sylvester-Kapoko/Receipts/printer"
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
		Payment: struct {
			Method string
			Amount decimal.Decimal
		}{
			Method: "Cash",
			Amount: decimal.RequireFromString("40.00"),
		},
		CreatedAt: time.Now(),
	}

	formatter := printer.NewPlainTextFormatter(32)   // if needed change back to 32
	receiptPrinter := printer.NewReceiptPrinter(formatter, os.Stdout)
	
                // print
               if err := receiptPrinter.Print(receipt); err != nil {
                   panic(err)
               }

}




