package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"

	"github.com/google/uuid"
	"github.com/herdifirdausss/belajar-vibe-coding/internal/models"
	"github.com/herdifirdausss/belajar-vibe-coding/internal/repository"
	"golang.org/x/crypto/argon2"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(email, password string) (*models.User, error) {
	// 1. Validate email
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, errors.New("invalid email format")
	}

	// 2. Validate password
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	// 3. Hash password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// 4. Generate UUID
	id := uuid.New()

	// 5. Create user struct
	user := &models.User{
		ID:       id,
		Email:    email,
		Password: hashedPassword,
	}

	// 6. Call repository
	return s.repo.CreateUser(user)
}

func hashPassword(password string) (string, error) {
	// Argon2id parameters
	time := uint32(1)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLen := uint32(32)
	saltLen := uint32(16)

	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	// Format: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, time, threads, b64Salt, b64Hash)
	return encoded, nil
}
