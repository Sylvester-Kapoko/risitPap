// internal/server/auth_handlers.go
package server

import (
    "html/template"
    "net/http"

    "github.com/Sylvester-Kapoko/risitPap/internal/store"
)

func HandleLogin(auth store.AuthStore) http.HandlerFunc {
    tmpl := template.Must(template.New("login").Parse(loginTemplate))
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            tmpl.Execute(w, nil)
            return
        }
        r.ParseForm()
        username := r.FormValue("username")
        password := r.FormValue("password")
        if _, err := auth.ValidateUser(username, password); err != nil {
            tmpl.Execute(w, map[string]string{"Error": "Invalid credentials."})
            return
        }
        setSessionCookie(w, username)
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func HandleLogout(auth store.AuthStore) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        http.SetCookie(w, &http.Cookie{
            Name:   "mk_session",
            MaxAge: -1,
            Path:   "/",
        })
        http.Redirect(w, r, "/login", http.StatusSeeOther)
    }
}