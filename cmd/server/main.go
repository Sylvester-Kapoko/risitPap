package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/Sylvester-Kapoko/risitPap/internal/server"
	"github.com/Sylvester-Kapoko/risitPap/internal/store"
	"github.com/Sylvester-Kapoko/risitPap/printer"
	"github.com/pkg/browser"
)

func main() {
	// --- License / Trial Check ---
	hasLicense := store.ValidLicense()
	if !hasLicense {
		daysLeft := store.TrialDaysLeft()
		if daysLeft <= 0 {
			fmt.Println("Your 30-day trial has expired.")
			fmt.Println("Whatsapp or call 0768592677 to purchase a license.")
			fmt.Println("Press Enter to exit.  .  .")
			_, _ = fmt.Scanln()
			os.Exit(1)
		}
		fmt.Printf("Trial mode: %d days remaining.\n", daysLeft)
	}

	// --- Database ---
	db, err := store.Open("receipts.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// --- Default admin user (if none exist) ---
	if hasUsers, err := db.HasUsers(); err == nil && !hasUsers {
		if err := db.CreateUser("admin", "admin"); err != nil {
			log.Printf("WARNING: could not create default admin user: %v", err)
		} else {
			log.Println("Default admin user created (admin / admin). Change password after first login.")
		}
	}

	formatter := printer.NewHtmlFormatter(printer.HtmlConfig{
		DeveloperName:  "Sylvester Kapoko",
		DeveloperPhone: "0768592677",
		ShowFooter:     true,
	})

	// --- Public routes (no authentication required) ---
	http.HandleFunc("/login", server.HandleLogin(db))
	http.HandleFunc("/logout", server.HandleLogout(db))
	http.HandleFunc("/register", server.HandleRegisterForm(db))
	http.HandleFunc("/register/save", server.HandleRegisterSave(db))

	// --- Protected routes (require login) ---
	http.HandleFunc("/", server.RequireLogin(server.HandleIndex(formatter, db)))
	http.HandleFunc("/print", server.RequireLogin(server.HandlePrint(formatter, db)))
	http.HandleFunc("/history", server.RequireLogin(server.HandleHistory(formatter, db)))
	http.HandleFunc("/view", server.RequireLogin(server.HandleView(formatter, db)))
	http.HandleFunc("/suggest", server.RequireLogin(server.HandleSuggest(db)))
	http.HandleFunc("/verify", server.RequireLogin(server.HandleVerify(db)))
	http.HandleFunc("/sync-kra", server.RequireLogin(server.HandleSyncKRA(db)))
	// NEW

	fmt.Println("Receipt Printer running at:")
	fmt.Println("  http://localhost:8080")
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			fmt.Printf("  http://%s:8080\n", ipnet.IP)
		}
	}

	if err := browser.OpenURL("http://localhost:8080"); err != nil {
		fmt.Println("could not open browser:", err)
	}
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
