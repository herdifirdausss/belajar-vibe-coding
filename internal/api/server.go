package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/herdifirdausss/belajar-vibe-coding/internal/handler"
	"github.com/herdifirdausss/belajar-vibe-coding/internal/repository"
	"github.com/herdifirdausss/belajar-vibe-coding/internal/service"
)

type Server struct {
	DB          *sql.DB
	UserHandler *handler.UserHandler
}

func NewServer(db *sql.DB) *Server {
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	userSvc := service.NewUserService(userRepo, sessionRepo)
	userHandler := handler.NewUserHandler(userSvc)

	return &Server{
		DB:          db,
		UserHandler: userHandler,
	}
}

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/users", s.UserHandler.RegisterUserHandler)
	mux.HandleFunc("/api/users/login", s.UserHandler.LoginHandler)
	mux.HandleFunc("/api/users/me", s.AuthMiddleware(s.UserHandler.GetMe))
	mux.HandleFunc("/api/users/logout", s.AuthMiddleware(s.UserHandler.LogoutHandler))
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := "UP"
	dbStatus := "UP"

	err := s.DB.Ping()
	if err != nil {
		dbStatus = "DOWN"
		status = "DEGRADED"
	}

	response := map[string]string{
		"status":   status,
		"database": dbStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
