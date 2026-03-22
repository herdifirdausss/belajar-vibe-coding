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
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password, created_at
		FROM users
		WHERE email = $1
	`
	user := &models.User{}
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	return user, nil
}

func (r *UserRepository) FindByToken(token string) (*models.User, error) {
	query := `
		SELECT u.id, u.email, u.created_at
		FROM users u
		INNER JOIN sessions s ON u.id = s.user_id
		WHERE s.token = $1
	`
	user := &models.User{}
	err := r.DB.QueryRow(query, token).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user by token: %w", err)
	}

	return user, nil
}
