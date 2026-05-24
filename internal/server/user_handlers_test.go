package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Sylvester-Kapoko/risitPap/domain"
)

func TestHandleChangePassword_Success(t *testing.T) {
	st := &mockStore{changePwdErr: nil}
	// We need a valid session cookie for getSessionUser to work.
	// Use the helper from auth_test.go: signTestCookie("admin")
	cookie := signTestCookie("admin")

	handler := HandleChangePassword(st)
	form := url.Values{}
	form.Set("old_password", "old")
	form.Set("new_password", "new")
	req := httptest.NewRequest("POST", "/changepassword", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/" {
		t.Errorf("expected redirect to /, got %s", loc)
	}
}

func TestHandleChangePassword_Failure(t *testing.T) {
	st := &mockStore{changePwdErr: errors.New("bad old password")}
	cookie := signTestCookie("admin")

	handler := HandleChangePassword(st)
	form := url.Values{}
	form.Set("old_password", "wrong")
	form.Set("new_password", "new")
	req := httptest.NewRequest("POST", "/changepassword", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Failed") {
		t.Error("expected error message on failure")
	}
}

func TestHandleListUsers_Success(t *testing.T) {
	st := &mockStore{
		usersList: []domain.User{
			{Username: "admin"},
			{Username: "staff1"},
		},
	}
	cookie := signTestCookie("admin")

	handler := HandleListUsers(st)
	req := httptest.NewRequest("GET", "/admin/users", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "admin") || !strings.Contains(body, "staff1") {
		t.Error("user list missing expected usernames")
	}
}

func TestHandleListUsers_Error(t *testing.T) {
	st := &mockStore{usersErr: errors.New("db error")}
	cookie := signTestCookie("admin")

	handler := HandleListUsers(st)
	req := httptest.NewRequest("GET", "/admin/users", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
}

func TestHandleAddUser_Success(t *testing.T) {
	st := &mockStore{users: make(map[string]string)} // needed for CreateUser mock
	cookie := signTestCookie("admin")

	handler := HandleAddUser(st)
	form := url.Values{}
	form.Set("username", "newuser")
	form.Set("password", "secret")
	req := httptest.NewRequest("POST", "/admin/users/add", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
}

func TestHandleAddUser_Duplicate(t *testing.T) {
	st := &mockStore{
		users: map[string]string{"existing": "hash"},
		createUserFn: func(username, password string) error {
			if username == "existing" {
				return errors.New("duplicate")
			}
			return nil
		},
	}

	cookie := signTestCookie("admin")
	handler := HandleAddUser(st)
	form := url.Values{}
	form.Set("username", "existing")
	form.Set("password", "newpass")
	req := httptest.NewRequest("POST", "/admin/users/add", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHandleDeleteUser_Success(t *testing.T) {
	st := &mockStore{deleteUserErr: nil}
	cookie := signTestCookie("admin")

	handler := HandleDeleteUser(st)
	req := httptest.NewRequest("GET", "/admin/users/delete?username=staff1", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
}

func TestHandleDeleteUser_LastUser(t *testing.T) {
	st := &mockStore{deleteUserErr: errors.New("cannot delete the last user")}
	cookie := signTestCookie("admin")

	handler := HandleDeleteUser(st)
	req := httptest.NewRequest("GET", "/admin/users/delete?username=admin", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}
