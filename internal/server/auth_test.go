// internal/server/auth_test.go
package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Sylvester-Kapoko/risitPap/domain"
)

// ---------- mock for AuthStore (now checks passwords) ----------
type mockAuthStore struct {
	users    map[string]string // username -> password (plain text for testing)
	hasUsers bool
}

func (m *mockAuthStore) HasUsers() (bool, error) { return m.hasUsers, nil }

func (m *mockAuthStore) CreateUser(username, password string) error {
	m.users[username] = password
	return nil
}

func (m *mockAuthStore) ValidateUser(username, password string) (*domain.User, error) {
	stored, ok := m.users[username]
	if !ok || stored != password {
		return nil, errors.New("invalid credentials")
	}
	return &domain.User{Username: username}, nil
}

// ---------- helper: create a valid signed cookie (uses production sessionSecret) ----------
func signTestCookie(username string) *http.Cookie {
	expiry := time.Now().Add(time.Hour).Format(time.RFC3339)
	value := username + "|" + expiry
	mac := hmac.New(sha256.New, sessionSecret)
	mac.Write([]byte(value))
	sig := hex.EncodeToString(mac.Sum(nil))[:16]
	return &http.Cookie{
		Name:  "mk_session",
		Value: value + "|" + sig,
	}
}

// ---------- tests ----------
func TestHandleLogin_GET(t *testing.T) {
	auth := &mockAuthStore{users: make(map[string]string), hasUsers: true}
	handler := HandleLogin(auth)
	req := httptest.NewRequest("GET", "/login", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<form") {
		t.Errorf("login page should contain a form")
	}
}

func TestHandleLogin_POST_Valid(t *testing.T) {
	auth := &mockAuthStore{users: map[string]string{"admin": "admin"}, hasUsers: true}
	handler := HandleLogin(auth)

	form := url.Values{}
	form.Set("username", "admin")
	form.Set("password", "admin")
	req := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
	// check for session cookie
	var session *http.Cookie
	for _, c := range rec.Result().Cookies() {
		if c.Name == "mk_session" {
			session = c
			break
		}
	}
	if session == nil {
		t.Fatal("expected session cookie mk_session")
	}
	if loc := rec.Header().Get("Location"); loc != "/" {
		t.Errorf("expected redirect to /, got %s", loc)
	}
}

func TestHandleLogin_POST_Invalid(t *testing.T) {
	auth := &mockAuthStore{users: map[string]string{"admin": "admin"}, hasUsers: true}
	handler := HandleLogin(auth)

	form := url.Values{}
	form.Set("username", "admin")
	form.Set("password", "wrong")
	req := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Invalid credentials") {
		t.Errorf("expected error message 'Invalid credentials'")
	}
}

func TestHandleLogout(t *testing.T) {
	auth := &mockAuthStore{}
	handler := HandleLogout(auth)
	req := httptest.NewRequest("GET", "/logout", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
	// cookie should be cleared
	for _, c := range rec.Result().Cookies() {
		if c.Name == "mk_session" && c.MaxAge == -1 {
			return // pass
		}
	}
	t.Error("expected session cookie to be cleared")
}

func TestRequireLogin_NoCookie(t *testing.T) {
	handler := RequireLogin(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/login" {
		t.Errorf("expected redirect to /login, got %s", loc)
	}
}

func TestRequireLogin_ValidSession(t *testing.T) {
	// create a valid cookie using the production session secret
	cookie := signTestCookie("admin")
	handler := RequireLogin(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}