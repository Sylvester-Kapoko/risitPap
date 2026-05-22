package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
  "log"
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

	db, err := store.Open("receipts.db")
	if err != nil {
		panic(err)
	}

	formatter := printer.NewHtmlFormatter(printer.HtmlConfig{
		DeveloperName:  "Sylvester Kapoko",
		DeveloperPhone: "0768592677",
		ShowFooter:     true,
	})
	http.HandleFunc("/", server.HandleIndex(formatter, db))
	http.HandleFunc("/print", server.HandlePrint(formatter, db))
	http.HandleFunc("/history", server.HandleHistory(formatter, db))
	http.HandleFunc("/view", server.HandleView(formatter, db))
	http.HandleFunc("/register", server.HandleRegisterForm(db))
	http.HandleFunc("/register/save", server.HandleRegisterSave(db))
	http.HandleFunc("/suggest", server.HandleSuggest(db))

	fmt.Println("Receipt Printer running at:")
	fmt.Println("  http://localhost:8080")
	// Print local IP for phone access
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
