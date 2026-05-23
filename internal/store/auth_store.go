// internal/store/auth_store.go
package store

import "github.com/Sylvester-Kapoko/risitPap/domain"

type AuthStore interface {
	HasUsers() (bool, error)
	CreateUser(username, password string) error
	ValidateUser(username, password string) (*domain.User, error)
}
