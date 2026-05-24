package server

import (
	"net/http"

	"github.com/Sylvester-Kapoko/risitPap/internal/store"
	"github.com/Sylvester-Kapoko/risitPap/printer"
)

func HandlePlainReceipt(st store.ReceiptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		receipt, err := st.Get(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		pf := printer.NewPlainTextFormatter(32) // 32 chars wide – typical thermal roll
		plain, err := pf.Format(receipt)
		if err != nil {
			http.Error(w, "Failed to format plain receipt", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(plain)) // #nosec G705
	}
}
