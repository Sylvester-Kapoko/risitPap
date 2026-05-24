// internal/store/interface.go
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

	// --- eTIMS audit log ---
	LogAction(username, action, detail string) error
	GetUnsyncedReceipts() ([]domain.Receipt, error)
	UpdateReceiptSync(r domain.Receipt) error
	CreateUser(username, password string) error
	ListUsers() ([]domain.User, error)
	DeleteUser(username string) error
	ChangePassword(username, oldPassword, newPassword string) error
}
