// internal/server/verify_handler.go
package server

import (
	"fmt"
	"net/http"

	"github.com/Sylvester-Kapoko/risitPap/internal/store"
)

func HandleVerify(st store.ReceiptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		receipt, err := st.Get(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		valid := receipt.Verify(signatureSecret)
		fmt.Fprintf(w, `<html><body><h1>Verification: %v</h1></body></html>`, valid)
	}
}
