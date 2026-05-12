package printer

import (
	"github.com/Sylvester-Kapoko/Receipts/domain"
)

type ReceiptFormatter interface {
	Format(r *domain.Receipt) string
}
