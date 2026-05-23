package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Sylvester-Kapoko/risitPap/domain"
	"github.com/Sylvester-Kapoko/risitPap/internal/store"
	"github.com/Sylvester-Kapoko/risitPap/printer"
	"github.com/shopspring/decimal"
)

var eatLocation = time.FixedZone("EAT", 3*60*60)

func nowInEAT() time.Time {
	return time.Now().UTC().In(eatLocation)
}

// getUsername safely reads the username from the signed session cookie.
// It returns an empty string if no valid session is present.
func getUsername(r *http.Request) string {
	cookie, err := r.Cookie("mk_session")
	if err != nil {
		return ""
	}
	parts := strings.Split(cookie.Value, "|")
	if len(parts) < 3 {
		return ""
	}
	return parts[0]
}

func HandleIndex(formatter *printer.HtmlFormatter, st store.ReceiptStore) http.HandlerFunc {
	tmpl := template.Must(template.New("index").Parse(indexHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		if !st.ValidLicense() && st.TrialDaysLeft() <= 0 {
			_, _ = w.Write([]byte(`<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Trial Expired</title>
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
		_ = tmpl.Execute(w, cfg)
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

		// --- Date handling (custom or auto) ---
		receiptDateStr := strings.TrimSpace(r.FormValue("receiptDate"))
		if receiptDateStr != "" {
			parsed, err := time.ParseInLocation("2006-01-02", receiptDateStr, eatLocation)
			if err != nil {
				receipt.CreatedAt = nowInEAT()
			} else {
				now := nowInEAT()
				receipt.CreatedAt = time.Date(parsed.Year(), parsed.Month(), parsed.Day(),
					now.Hour(), now.Minute(), now.Second(), 0, eatLocation)
			}
		} else {
			receipt.CreatedAt = nowInEAT()
		}

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

		// -------- eTIMS digital signature --------
		receipt.Sign(signatureSecret)

		if err := st.Save(receipt); err != nil {
			http.Error(w, "Failed to save receipt", http.StatusInternalServerError)
			return
		}

		// -------- Audit log (print) --------
		if user := getUsername(r); user != "" {
			_ = st.LogAction(user, "print", receipt.TransactionID)
		}

		formatted, err := formatter.Format(receipt)
		if err != nil {
			http.Error(w, "Failed to format receipt", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(formatted))
	}
}

func HandleRegisterForm(st store.ReceiptStore) http.HandlerFunc {
	tmpl := template.Must(template.New("register").Parse(registerHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		cfg, _ := st.LoadConfig()
		if cfg == nil {
			cfg = &domain.StoreConfig{}
		}
		_ = tmpl.Execute(w, cfg)
	}
}

func HandleRegisterSave(st store.ReceiptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("parseform: %v", err)
		}
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

		// -------- Audit log (config change) --------
		if user := getUsername(r); user != "" {
			_ = st.LogAction(user, "config_change", "")
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func generateTransactionID() string {
	now := nowInEAT()
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		log.Printf("rand.Read: %v", err)
	}
	return fmt.Sprintf("INV-%s-%s", now.Format("20060102"), hex.EncodeToString(b))
}

func parseForm(r *http.Request) *domain.Receipt {
	if err := r.ParseForm(); err != nil {
		log.Printf("parseform: %v", err)
	}
	taxRate, _ := decimal.NewFromString(r.FormValue("taxRate"))
	payAmt, _ := decimal.NewFromString(r.FormValue("paymentAmount"))
	itemsStr := r.FormValue("items")
	var items []domain.ReceiptItem
	for _, part := range strings.Split(itemsStr, ";") {
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
		TransactionID: r.FormValue("transactionID"),
		StoreName:     r.FormValue("storeName"),
		StoreAddr:     r.FormValue("storeAddr"),
		TaxRate:       taxRate,
		Items:         items,
		Payment: domain.Payment{
			Method: r.FormValue("paymentMethod"),
			Amount: payAmt,
		},
		CreatedAt:  nowInEAT(),
		SyncStatus: "pending", // ← new line, correctly inside the struct
	}
}