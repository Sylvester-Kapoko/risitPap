package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Sylvester-Kapoko/risitPap/domain"
)

// --------------------------------------------------------------------
//  Mock formatter
// --------------------------------------------------------------------

type mockFormatter struct {
	output string
	err    error
}

func (m *mockFormatter) Format(r *domain.Receipt) (string, error) {
	return m.output, m.err
}

// --------------------------------------------------------------------
//  Mock store
// --------------------------------------------------------------------

type mockStore struct {
	receipts    map[string]*domain.Receipt
	config      *domain.StoreConfig
	license     bool
	trialDays   int
	suggestFn   func(prefix string) ([]map[string]string, error)
	unsynced    []domain.Receipt
	unsyncedErr error
	updateErr   error
}

func (m *mockStore) Save(r *domain.Receipt) error {
	if m.receipts == nil {
		m.receipts = make(map[string]*domain.Receipt)
	}
	m.receipts[r.TransactionID] = r
	return nil
}

func (m *mockStore) Get(id string) (*domain.Receipt, error) {
	if m.receipts == nil {
		return nil, errors.New("not found")
	}
	r, ok := m.receipts[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return r, nil
}

func (m *mockStore) List(since time.Time) ([]domain.Receipt, error) {
	var all []domain.Receipt
	for _, r := range m.receipts {
		all = append(all, *r)
	}
	return all, nil
}

func (m *mockStore) LoadConfig() (*domain.StoreConfig, error) {
	return m.config, nil
}

func (m *mockStore) SaveConfig(cfg *domain.StoreConfig) error {
	m.config = cfg
	return nil
}

func (m *mockStore) ValidLicense() bool {
	return m.license
}

func (m *mockStore) TrialDaysLeft() int {
	return m.trialDays
}

func (m *mockStore) SuggestItems(prefix string) ([]map[string]string, error) {
	if m.suggestFn != nil {
		return m.suggestFn(prefix)
	}
	return nil, nil
}

// unsynced is a slice of receipts that GetUnsyncedReceipts will return.
func (m *mockStore) GetUnsyncedReceipts() ([]domain.Receipt, error) {
	return m.unsynced, m.unsyncedErr
}

// UpdateReceiptSync is a no-op that records the call.
func (m *mockStore) UpdateReceiptSync(r domain.Receipt) error {
	return m.updateErr
}

func (m *mockStore) LogAction(username, action, detail string) error {
	return nil
}

// --------------------------------------------------------------------
//  Date handling tests
// --------------------------------------------------------------------

func TestHandlePrint_DateHandling(t *testing.T) {
	// ----------------------------------------------------------------
	// Test 1: Custom date is parsed correctly in EAT
	// ----------------------------------------------------------------
	st := &mockStore{
		license:  true,
		receipts: make(map[string]*domain.Receipt),
		config:   &domain.StoreConfig{Currency: "Ksh"},
	}
	fm := &mockFormatter{output: "<p>ok</p>"}

	handler := HandlePrint(fm, st)

	body := "storeName=Test&storeAddr=Addr&taxRate=0.16&paymentMethod=Cash&paymentAmount=100&items=Item,1,50.00&receiptDate=2025-05-15"
	req := httptest.NewRequest("POST", "/print", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("test 1: expected 200, got %d", rec.Code)
	}

	var saved *domain.Receipt
	for _, r := range st.receipts {
		saved = r
		break
	}
	if saved == nil {
		t.Fatal("test 1: receipt was not saved")
	}

	if saved.CreatedAt.Year() != 2025 || saved.CreatedAt.Month() != time.May || saved.CreatedAt.Day() != 15 {
		t.Errorf("test 1: expected date 2025-05-15, got %s", saved.CreatedAt.Format("2006-01-02"))
	}

	if saved.CreatedAt.Location().String() != "EAT" {
		t.Errorf("test 1: expected location EAT, got %s", saved.CreatedAt.Location().String())
	}

	// ----------------------------------------------------------------
	// Test 2: Auto date (no receiptDate) is close to now in EAT
	// ----------------------------------------------------------------
	st2 := &mockStore{
		license:  true,
		receipts: make(map[string]*domain.Receipt),
		config:   &domain.StoreConfig{Currency: "Ksh"},
	}
	handler2 := HandlePrint(fm, st2)

	body2 := "storeName=Test&storeAddr=Addr&taxRate=0.16&paymentMethod=Cash&paymentAmount=100&items=Item,1,50.00"
	req2 := httptest.NewRequest("POST", "/print", strings.NewReader(body2))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec2 := httptest.NewRecorder()
	handler2.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("test 2: expected 200, got %d", rec2.Code)
	}

	var saved2 *domain.Receipt
	for _, r := range st2.receipts {
		saved2 = r
		break
	}
	if saved2 == nil {
		t.Fatal("test 2: receipt was not saved")
	}

	// Compare against time.Now().In(eatLocation) — matches the fixed handler
	nowEAT := time.Now().In(eatLocation)
	diff := saved2.CreatedAt.Sub(nowEAT)
	if diff < -2*time.Second || diff > 2*time.Second {
		t.Errorf("test 2: auto date should be close to now EAT, got %s, expected near %s",
			saved2.CreatedAt.Format("2006-01-02 15:04:05"),
			nowEAT.Format("2006-01-02 15:04:05"))
	}

	if saved2.CreatedAt.Location().String() != "EAT" {
		t.Errorf("test 2: expected location EAT, got %s", saved2.CreatedAt.Location().String())
	}
}
