package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sylvester-Kapoko/risitPap/domain"
)

func TestHandleSyncKRA_Success(t *testing.T) {
	st := &mockStore{
		receipts: make(map[string]*domain.Receipt),
		unsynced: []domain.Receipt{
			{TransactionID: "INV-001", SyncStatus: "pending"},
			{TransactionID: "INV-002", SyncStatus: "pending"},
		},
	}
	handler := HandleSyncKRA(st)

	req := httptest.NewRequest("GET", "/sync-kra", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "2 receipts synced") {
		t.Errorf("expected '2 receipts synced', got %q", body)
	}

	// Verify the receipts in the mock's unsynced slice were modified
	if st.unsynced[0].FDN != "FDN-DUMMY-INV-001" {
		t.Errorf("FDN[0] = %q, want %q", st.unsynced[0].FDN, "FDN-DUMMY-INV-001")
	}
	if st.unsynced[0].AntiFakeCode != "AF-DUMMY-INV-001" {
		t.Errorf("AntiFakeCode[0] = %q, want %q", st.unsynced[0].AntiFakeCode, "AF-DUMMY-INV-001")
	}
	if st.unsynced[0].SyncStatus != "synced" {
		t.Errorf("SyncStatus[0] = %q, want synced", st.unsynced[0].SyncStatus)
	}
}

func TestHandleSyncKRA_EmptyList(t *testing.T) {
	st := &mockStore{
		unsynced: []domain.Receipt{},
	}
	handler := HandleSyncKRA(st)

	req := httptest.NewRequest("GET", "/sync-kra", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "0 receipts synced") {
		t.Errorf("expected '0 receipts synced', got %q", body)
	}
}

func TestHandleSyncKRA_StoreError(t *testing.T) {
	st := &mockStore{
		unsyncedErr: errors.New("database offline"),
	}
	handler := HandleSyncKRA(st)

	req := httptest.NewRequest("GET", "/sync-kra", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
}
