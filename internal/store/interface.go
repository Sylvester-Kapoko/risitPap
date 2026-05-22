package store

import (
	"time"

	"github.com/Sylvester-Kapoko/risitPap/domain"
)

type ReceiptStore interface {
	Save(r *domain.Receipt) error
	Get(id string) (*domain.Receipt, error)
	List(since time.Time) ([]domain.Receipt, error)
	LoadConfig() (*domain.StoreConfig, error)
	SaveConfig(cfg *domain.StoreConfig) error
	ValidLicense() bool
	TrialDaysLeft() int
	SuggestItems(prefix string) ([]map[string]string, error)
}
