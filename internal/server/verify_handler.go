// internal/server/verify_handler.go
package server

import (
	"html/template"
	"net/http"

	"github.com/Sylvester-Kapoko/risitPap/internal/store"
)

var verifyTmpl = template.Must(template.New("verify").Parse(`<html><body><h1>Verification: {{.}}</h1></body></html>`))

func HandleVerify(st store.ReceiptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		receipt, err := st.Get(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		valid := receipt.Verify(signatureSecret)
		_ = verifyTmpl.Execute(w, valid)
	}
}
