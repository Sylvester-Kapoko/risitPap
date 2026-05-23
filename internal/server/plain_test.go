package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Sylvester-Kapoko/risitPap/domain"
	"github.com/shopspring/decimal"
)

func TestHandlePlainReceipt_Success(t *testing.T) {
	// Create a receipt with known values
	r := &domain.Receipt{
		TransactionID: "INV-001",
		StoreName:     "Test Store",
		StoreAddr:     "123 Main St",
		Currency:      "Ksh",
		Items: []domain.ReceiptItem{
			{Name: "Widget", Qty: 2, UnitPrice: decimal.NewFromFloat(10.50)},
		},
		TaxRate:   decimal.Zero,
		Payment:   domain.Payment{Method: "Cash", Amount: decimal.NewFromFloat(21.00)},
		CreatedAt: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	st := &mockStore{
		receipts: map[string]*domain.Receipt{"INV-001": r},
	}

	handler := HandlePlainReceipt(st)

	req := httptest.NewRequest("GET", "/plain?id=INV-001", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		t.Errorf("expected Content-Type text/plain, got %q", contentType)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Test Store") {
		t.Error("plain output missing store name")
	}
	if !strings.Contains(body, "Widget") {
		t.Error("plain output missing item name")
	}
}

func TestHandlePlainReceipt_NotFound(t *testing.T) {
	st := &mockStore{
		receipts: make(map[string]*domain.Receipt),
	}
	handler := HandlePlainReceipt(st)

	req := httptest.NewRequest("GET", "/plain?id=nonexistent", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}
