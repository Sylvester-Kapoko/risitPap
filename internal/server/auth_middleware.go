// internal/server/auth_middleware.go
package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	sessionCookieName = "mk_session"
	sessionDuration   = 24 * time.Hour
	cookieParts       = 3
)

var sessionSecret = loadSecret()

func loadSecret() []byte {
	s := os.Getenv("SESSION_SECRET")
	if s == "" {
		log.Println("warn: SESSION_SECRET not set, using insecure default")
		// #nosec G101 – fallback for development only
		return []byte("change-me-to-32-random-bytes")
	}
	if len(s) < 32 {
		log.Println("warn: SESSION_SECRET should be at least 32 bytes")
	}
	return []byte(s)
}

// sessionPayload holds the decoded fields of a session cookie value.
type sessionPayload struct {
	Username  string
	ExpiryStr string
	Sig       string
}

var errInvalidCookie = errors.New("invalid session cookie")

func setSessionCookie(w http.ResponseWriter, username string) {
	expiry := time.Now().Add(sessionDuration)
	expiryStr := expiry.Format(time.RFC3339)
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    username + "|" + expiryStr + "|" + signPayload(username, expiryStr),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expiry,
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

// RequireLogin redirects to /login if the request carries no valid session.
func RequireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, _, ok := GetSessionUser(r); !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// GetSessionUser extracts the username from the signed session cookie.
// Returns ("", "", false) if the cookie is missing, expired, or tampered with.
func GetSessionUser(r *http.Request) (username, role string, ok bool) {
	p, err := parseSessionCookie(r)
	if err != nil {
		return "", "", false
	}
	expiry, err := time.Parse(time.RFC3339, p.ExpiryStr)
	if err != nil || time.Now().After(expiry) {
		return "", "", false
	}
	if !hmac.Equal([]byte(p.Sig), []byte(signPayload(p.Username, p.ExpiryStr))) {
		return "", "", false
	}
	return p.Username, "", true
}

// parseSessionCookie reads and splits the raw cookie value.
func parseSessionCookie(r *http.Request) (*sessionPayload, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(cookie.Value, "|", cookieParts)
	if len(parts) != cookieParts {
		return nil, errInvalidCookie
	}
	return &sessionPayload{
		Username:  parts[0],
		ExpiryStr: parts[1],
		Sig:       parts[2],
	}, nil
}

// signPayload returns the full hex HMAC-SHA256 of "username|expiryStr".
func signPayload(username, expiryStr string) string {
	mac := hmac.New(sha256.New, sessionSecret)
	mac.Write([]byte(username + "|" + expiryStr))
	return hex.EncodeToString(mac.Sum(nil))
}
