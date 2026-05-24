package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Sylvester-Kapoko/risitPap/internal/store"
)

func HandleLogin(auth store.AuthStore) http.HandlerFunc {
	tmpl := template.Must(template.New("login").Parse(loginTemplate))
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			_ = tmpl.Execute(w, nil) // write error ignored (non-fatal)
			return
		}
		if err := r.ParseForm(); err != nil {
			log.Printf("parseform: %v", err)
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		if _, err := auth.ValidateUser(username, password); err != nil {
			_ = tmpl.Execute(w, map[string]string{"Error": "Invalid credentials."})
			return
		}
		setSessionCookie(w, username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
func HandleLogout(auth store.AuthStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clearSessionCookie(w)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
