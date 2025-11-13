package users

import (
	"context"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	err := h.service.RegisterUser(context.Background())
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
	}
}
