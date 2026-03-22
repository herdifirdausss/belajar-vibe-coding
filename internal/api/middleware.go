package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/herdifirdausss/belajar-vibe-coding/internal/models"
)

func (s *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "unauthorized"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "unauthorized"})
			return
		}

		token := parts[1]
		ctx := context.WithValue(r.Context(), models.TokenKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
