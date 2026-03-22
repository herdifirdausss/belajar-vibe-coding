package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/herdifirdausss/belajar-vibe-coding/internal/models"
	"github.com/herdifirdausss/belajar-vibe-coding/internal/repository"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidEmailFormat     = errors.New("invalid email format")
	ErrEmailTooLong           = errors.New("email must not exceed 255 characters")
	ErrPasswordTooShort       = errors.New("password must be at least 8 characters long")
	ErrEmailOrPasswordMissing = errors.New("email and password are required")
	ErrInvalidCredentials     = errors.New("email or password incorrect")
	ErrUnauthorized           = errors.New("unauthorized")
)

type UserService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

func NewUserService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *UserService {
	return &UserService{userRepo: userRepo, sessionRepo: sessionRepo}
}

func (s *UserService) Register(email, password string) (*models.User, error) {
	// 1. Validate email
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, ErrInvalidEmailFormat
	}

	// 1.b Validate email length
	if len(email) > 255 {
		return nil, ErrEmailTooLong
	}

	// 2. Validate password
	if len(password) < 8 {
		return nil, ErrPasswordTooShort
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
	return s.userRepo.CreateUser(user)
}

func (s *UserService) Login(email, password string) (*string, error) {
	// 1. Validate input
	if email == "" || password == "" {
		return nil, ErrEmailOrPasswordMissing
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, ErrInvalidEmailFormat
	}

	// 2. Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// 3. Compare password
	match, err := comparePassword(user.Password, password)
	if err != nil {
		return nil, fmt.Errorf("error comparing password: %w", err)
	}

	if !match {
		return nil, ErrInvalidCredentials
	}

	// 4. Generate token
	token := uuid.New().String()

	// 5. Create session
	session := &models.Session{
		ID:     uuid.New(),
		Token:  token,
		UserID: user.ID,
	}

	// 6. Save session
	err = s.sessionRepo.Create(session)
	if err != nil {
		return nil, fmt.Errorf("error creating session: %w", err)
	}

	return &token, nil
}

func (s *UserService) Me(token string) (*models.User, error) {
	if token == "" {
		return nil, ErrUnauthorized
	}

	user, err := s.userRepo.FindByToken(token)
	if err != nil {
		return nil, fmt.Errorf("error finding current user: %w", err)
	}

	if user == nil {
		return nil, ErrUnauthorized
	}

	return user, nil
}

func (s *UserService) Logout(token string) error {
	if token == "" {
		return ErrUnauthorized
	}

	return s.sessionRepo.DeleteByToken(token)
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

func comparePassword(encodedHash, password string) (bool, error) {
	// Format: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, fmt.Errorf("error parsing version: %w", err)
	}

	var memory, time uint32
	var threads uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, fmt.Errorf("error parsing parameters: %w", err)
	}

	b64Salt := parts[4]
	b64Hash := parts[5]

	salt, err := base64.RawStdEncoding.DecodeString(b64Salt)
	if err != nil {
		return false, fmt.Errorf("error decoding salt: %w", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(b64Hash)
	if err != nil {
		return false, fmt.Errorf("error decoding hash: %w", err)
	}

	otherHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(hash)))

	return string(hash) == string(otherHash), nil
}
