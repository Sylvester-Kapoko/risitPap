// internal/server/auth_middleware.go
package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

var sessionSecret = []byte("change-me-to-32-random-bytes")

func setSessionCookie(w http.ResponseWriter, username string) {
	expiry := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	value := username + "|" + expiry
	mac := hmac.New(sha256.New, sessionSecret)
	mac.Write([]byte(value))
	sig := hex.EncodeToString(mac.Sum(nil))[:16]
	http.SetCookie(w, &http.Cookie{
		Name:     "mk_session",
		Value:    value + "|" + sig,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}

func RequireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("mk_session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		parts := strings.Split(cookie.Value, "|")
		if len(parts) != 3 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		username, expiryStr, sig := parts[0], parts[1], parts[2]
		expiry, err := time.Parse(time.RFC3339, expiryStr)
		if err != nil || time.Now().After(expiry) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		mac := hmac.New(sha256.New, sessionSecret)
		mac.Write([]byte(username + "|" + expiryStr))
		expected := hex.EncodeToString(mac.Sum(nil))[:16]
		if sig != expected {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
