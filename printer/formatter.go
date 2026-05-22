package printer

import "github.com/Sylvester-Kapoko/risitPap/domain"

type ReceiptFormatter interface {
	Format(r *domain.Receipt) (string, error)
}
