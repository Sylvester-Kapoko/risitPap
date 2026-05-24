package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Sylvester-Kapoko/risitPap/internal/store"
)

func HandleChangePassword(st store.ReceiptStore) http.HandlerFunc {
	tmpl := template.Must(template.New("changepassword").Parse(changePasswordTemplate))
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			_ = tmpl.Execute(w, nil)
			return
		}
		if err := r.ParseForm(); err != nil {
			log.Printf("parseform: %v", err)
		}
		username, _, _ := GetSessionUser(r)
		oldPass := r.FormValue("old_password")
		newPass := r.FormValue("new_password")
		if err := st.ChangePassword(username, oldPass, newPass); err != nil {
			_ = tmpl.Execute(w, map[string]string{"Error": "Failed. Check your old password."})
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func HandleListUsers(st store.ReceiptStore) http.HandlerFunc {
	tmpl := template.Must(template.New("userlist").Parse(userListTemplate))
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := st.ListUsers()
		if err != nil {
			http.Error(w, "Could not list users", http.StatusInternalServerError)
			return
		}
		_, role, _ := GetSessionUser(r)
		_ = tmpl.Execute(w, map[string]interface{}{
			"Users":        users,
			"IsSupervisor": role == "supervisor",
		})
	}
}

func HandleAddUser(st store.ReceiptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("parseform: %v", err)
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		if err := st.CreateUser(username, password); err != nil {
			http.Error(w, "User already exists or invalid data", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
	}
}

func HandleDeleteUser(st store.ReceiptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		if err := st.DeleteUser(username); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
	}
}
