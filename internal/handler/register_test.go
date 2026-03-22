package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/herdifirdausss/belajar-vibe-coding/internal/service"
)

func TestUserHandler_RegisterUserHandler_EmailLength(t *testing.T) {
	// We need to pass a service that will return the error we want.
	// Since we are testing the handler's mapping of the error from the service,
	// we can use a real service but with nil repositories, because the validation happens BEFORE repo calls.
	svc := service.NewUserService(nil, nil)
	h := NewUserHandler(svc)

	// Test case: email 256 characters
	email256 := strings.Repeat("a", 247) + "@test.com"
	reqBody, _ := json.Marshal(RegisterRequest{
		Email:    email256,
		Password: "password123",
	})

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	h.RegisterUserHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error":"email must not exceed 255 characters"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
