package service

import (
	"errors"
	"strings"
	"testing"
)

func TestUserService_Register_EmailLength(t *testing.T) {
	svc := &UserService{}

	// Test case: email 256 characters
	// This should return an error BEFORE calling the repository.
	email256 := strings.Repeat("a", 247) + "@test.com"
	_, err := svc.Register(email256, "password123")
	if err == nil {
		t.Fatal("expected error for email length 256, but got nil")
	}
	if !errors.Is(err, ErrEmailTooLong) {
		t.Errorf("expected ErrEmailTooLong, but got: %v", err)
	}

	// Test case: invalid email format
	_, err = svc.Register("invalid-email", "password123")
	if err == nil {
		t.Fatal("expected error for invalid email format, but got nil")
	}
	if !errors.Is(err, ErrInvalidEmailFormat) {
		t.Errorf("expected ErrInvalidEmailFormat, but got: %v", err)
	}
}
