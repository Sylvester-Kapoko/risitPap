package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Sylvester-Kapoko/risitPap/domain"
	"github.com/Sylvester-Kapoko/risitPap/internal/store"
	"github.com/Sylvester-Kapoko/risitPap/printer"
	"github.com/shopspring/decimal"
)

// eatLocation is East Africa Time (UTC+3).
var eatLocation = time.FixedZone("EAT", 3*60*60)

// nowInEAT returns the current time in EAT.
func nowInEAT() time.Time {
	return time.Now().In(eatLocation)
}

func HandleIndex(formatter *printer.HtmlFormatter, st store.ReceiptStore) http.HandlerFunc {
	tmpl := template.Must(template.New("index").Parse(indexHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		if !st.ValidLicense() && st.TrialDaysLeft() <= 0 {
			w.Write([]byte(`<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Trial Expired</title>
            <style>body{font-family:system-ui,sans-serif;max-width:500px;margin:auto;padding:20px;text-align:center;}
            a{color:#007bff;}</style></head><body>
            <h1>⏰ Trial Expired</h1>
            <p>Your 30-day trial has ended.</p>
            <p>Your receipts are still available.</p>
            <p><a href="/history">📋 View Receipt History</a></p>
            <p>Purchase a license key to resume printing new receipts.</p>
            <p><a href="https://wa.me/254768592677" style="font-size:1.2em;font-weight:bold;">🛒 Buy Now (Ksh 2000)</a></p>
            </body></html>`))
			return
		}
		cfg, _ := st.LoadConfig()
		if cfg == nil {
			cfg = &domain.StoreConfig{}
		}
		tmpl.Execute(w, cfg)
	}
}

func HandlePrint(formatter printer.ReceiptFormatter, st store.ReceiptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if !st.ValidLicense() && st.TrialDaysLeft() <= 0 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		receipt := parseForm(r)
		receipt.TransactionID = generateTransactionID()
		receipt.CreatedAt = resolveReceiptDate(r.FormValue("receiptDate"))

		if cfg, err := st.LoadConfig(); err == nil && cfg != nil {
			receipt.StorePhone = cfg.StorePhone
			receipt.StoreTaxID = cfg.StoreTaxID
			receipt.Currency = cfg.Currency
			receipt.HasVAT = cfg.VATRegistered
			if !cfg.VATRegistered {
				receipt.TaxRate = decimal.Zero
			}
		}
		if receipt.Currency == "" {
			receipt.Currency = "Ksh"
		}

		if err := st.Save(receipt); err != nil {
			http.Error(w, "Failed to save receipt", http.StatusInternalServerError)
			return
		}

		formatted, err := formatter.Format(receipt)
		if err != nil {
			http.Error(w, "Failed to format receipt", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(formatted))
	}
}

// resolveReceiptDate parses a YYYY-MM-DD date string into EAT time,
// preserving the current EAT clock time. Falls back to now if blank or invalid.
func resolveReceiptDate(dateStr string) time.Time {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return nowInEAT()
	}
	parsed, err := time.ParseInLocation("2006-01-02", dateStr, eatLocation)
	if err != nil {
		return nowInEAT()
	}
	now := nowInEAT()
	return time.Date(
		parsed.Year(), parsed.Month(), parsed.Day(),
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		eatLocation,
	)
}

func HandleRegisterForm(st store.ReceiptStore) http.HandlerFunc {
	tmpl := template.Must(template.New("register").Parse(registerHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		cfg, _ := st.LoadConfig()
		if cfg == nil {
			cfg = &domain.StoreConfig{}
		}
		tmpl.Execute(w, cfg)
	}
}

func HandleRegisterSave(st store.ReceiptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		vatRegistered := r.FormValue("vatRegistered") == "true"
		cfg := &domain.StoreConfig{
			StoreName:     r.FormValue("storeName"),
			StoreAddr:     r.FormValue("storeAddr"),
			StorePhone:    r.FormValue("storePhone"),
			StoreTaxID:    r.FormValue("storeTaxID"),
			VATRegistered: vatRegistered,
			Currency:      r.FormValue("currency"),
			TemplateId:    r.FormValue("templateId"),
			LogoPath:      r.FormValue("logoPath"),
		}
		if err := st.SaveConfig(cfg); err != nil {
			http.Error(w, "Failed to save settings", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func generateTransactionID() string {
	b := make([]byte, 3)
	rand.Read(b)
	return fmt.Sprintf("INV-%s-%s", nowInEAT().Format("20060102"), hex.EncodeToString(b))
}

func parseForm(r *http.Request) *domain.Receipt {
	r.ParseForm()
	taxRate, _ := decimal.NewFromString(r.FormValue("taxRate"))
	payAmt, _ := decimal.NewFromString(r.FormValue("paymentAmount"))
	var items []domain.ReceiptItem
	for _, part := range strings.Split(r.FormValue("items"), ";") {
		if part == "" {
			continue
		}
		fields := strings.Split(part, ",")
		if len(fields) < 3 {
			continue
		}
		qty, _ := strconv.Atoi(strings.TrimSpace(fields[1]))
		price, _ := decimal.NewFromString(strings.TrimSpace(fields[2]))
		items = append(items, domain.ReceiptItem{
			Name:      strings.TrimSpace(fields[0]),
			Qty:       qty,
			UnitPrice: price,
		})
	}
	return &domain.Receipt{
		StoreName: r.FormValue("storeName"),
		StoreAddr: r.FormValue("storeAddr"),
		TaxRate:   taxRate,
		Items:     items,
		Payment: domain.Payment{
			Method: r.FormValue("paymentMethod"),
			Amount: payAmt,
		},
	}
}