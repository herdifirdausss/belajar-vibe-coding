package repository

import (
	"database/sql"
	"fmt"

	"github.com/herdifirdausss/belajar-vibe-coding/internal/models"
)

type SessionRepository struct {
	DB *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (r *SessionRepository) Create(session *models.Session) error {
	query := `
		INSERT INTO sessions (id, token, user_id)
		VALUES ($1, $2, $3)
		RETURNING created_at
	`
	err := r.DB.QueryRow(query, session.ID, session.Token, session.UserID).Scan(&session.CreatedAt)
	if err != nil {
		return fmt.Errorf("error inserting session: %w", err)
	}

	return nil
}
