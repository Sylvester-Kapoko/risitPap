package server

import (
	"encoding/base64"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/Sylvester-Kapoko/risitPap/domain"
	"github.com/Sylvester-Kapoko/risitPap/internal/store"
	"github.com/Sylvester-Kapoko/risitPap/printer"
)

const historyHTML = `<!DOCTYPE html>
<html><head><meta charset="UTF-8"><title>Receipt History</title>
<style>body{font-family:monospace;max-width:700px;margin:auto;padding:20px;}
table{width:100%;border-collapse:collapse;}td,th{padding:5px;text-align:left;}
a{color:blue;}</style></head><body>
<h1>Receipt History</h1><a href="/">New Receipt</a>
<table><tr><th>ID</th><th>Date</th><th>Store</th><th>Total</th><th>View</th></tr>
{{range .}}<tr><td>{{.TransactionID}}</td><td>{{.CreatedAt.Format "2006-01-02 15:04"}}</td>
<td>{{.StoreName}}</td><td>{{.Total}}</td>
<td><a href="/view?id={{.TransactionID}}">View</a></td></tr>{{end}}</table>
</body></html>`

func HandleHistory(formatter *printer.HtmlFormatter, st *store.Store) http.HandlerFunc {
	tmpl := template.Must(template.New("history").Parse(historyHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		receipts, _ := st.List(time.Time{})
		type row struct {
			TransactionID, StoreName, Total string
			CreatedAt                       time.Time
		}
		rows := []row{}
		for _, rec := range receipts {
			currency := rec.Currency
			if currency == "" {
				currency = "Ksh"
			}
			rows = append(rows, row{rec.TransactionID, rec.StoreName,
				domain.DisplayCurrency(rec.Total(), currency), rec.CreatedAt})
		}
		_ = tmpl.Execute(w, rows)
	}
}

// receiptWrapper is used to pass already-safe HTML to the template.
type receiptWrapper struct {
	HTML template.HTML
}

func HandleView(formatter *printer.HtmlFormatter, st store.ReceiptStore) http.HandlerFunc {
	tmpl := template.Must(template.New("receiptview").Parse(`<html><body>{{.HTML}}</body></html>`))
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		receipt, err := st.Get(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// Load the logo if a path is configured
		if cfg, err := st.LoadConfig(); err == nil && cfg != nil && cfg.LogoPath != "" {
			if imgData, err := os.ReadFile(cfg.LogoPath); err == nil {
				receipt.LogoBase64 = base64.StdEncoding.EncodeToString(imgData)
			}
		}

		formatted, err := formatter.Format(receipt)
		if err != nil {
			http.Error(w, "Failed to format receipt", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = tmpl.Execute(w, receiptWrapper{HTML: template.HTML(formatted)}) // #nosec G203
	}
}
