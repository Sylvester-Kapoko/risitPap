package server

import (
	"fmt"
	"net/http"

	"github.com/Sylvester-Kapoko/risitPap/internal/store"
)

func HandleSyncKRA(st store.ReceiptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		unsynced, err := st.GetUnsyncedReceipts()
		if err != nil {
			http.Error(w, "Failed to fetch unsynced receipts", http.StatusInternalServerError)
			return
		}

		count := 0
		for i := range unsynced {
			rcpt := &unsynced[i] // pointer to the actual slice element
			rcpt.FDN = "FDN-DUMMY-" + rcpt.TransactionID
			rcpt.AntiFakeCode = "AF-DUMMY-" + rcpt.TransactionID
			rcpt.SyncStatus = "synced"

			if err := st.UpdateReceiptSync(*rcpt); err != nil {
				continue
			}
			count++
		}

		fmt.Fprintf(w, "%d receipts synced (dummy)", count)
	}
}
