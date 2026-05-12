package handler

import (
	"fmt"
	"net/http"
                "github.com/Sylvester-Kapoko/Receipts/service"
)

type Handler struct {
	service *service.UserService
}

func NewHandler(s *service.UserService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := h.service.GetUser(1)
	fmt.Fprintln(w, user)
}
