package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/herdifirdausss/belajar-vibe-coding/internal/models"
	"github.com/herdifirdausss/belajar-vibe-coding/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	user, err := h.svc.Register(req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidEmailFormat) || errors.Is(err, service.ErrEmailTooLong) || errors.Is(err, service.ErrPasswordTooShort) {
			status = http.StatusBadRequest
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Ensure password is not in the JSON response
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidEmailFormat) || errors.Is(err, service.ErrEmailOrPasswordMissing) {
			status = http.StatusBadRequest
		} else if errors.Is(err, service.ErrInvalidCredentials) {
			status = http.StatusUnauthorized
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	response := models.LoginResponse{
		Data: models.LoginData{
			Token: *token,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, ok := r.Context().Value(models.TokenKey).(string)
	if !ok || token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "unauthorized"})
		return
	}

	user, err := h.svc.Me(token)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrUnauthorized) {
			status = http.StatusUnauthorized
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	resp := map[string]interface{}{
		"data": map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, ok := r.Context().Value(models.TokenKey).(string)
	if !ok || token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "unauthorized"})
		return
	}

	err := h.svc.Logout(token)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrUnauthorized) {
			status = http.StatusUnauthorized
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	resp := map[string]interface{}{
		"data": map[string]interface{}{
			"message": "success logout",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
