package repository

import (
	"database/sql"
	"fmt"

	"github.com/herdifirdausss/belajar-vibe-coding/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (id, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	err := r.DB.QueryRow(query, user.ID, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error inserting user: %w", err)
	}

	return user, nil
}
