package service

import (
 
            "github.com/Sylvester-Kapoko/Receipts/repository"
)


type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) GetUser(id int) string {
	// Add some business logic here
	return s.repo.FindUser(id)
}

