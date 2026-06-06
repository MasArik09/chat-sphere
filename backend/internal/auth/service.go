package auth

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"chatsphere/internal/users"
	"chatsphere/pkg/config"
)

var (
	// ErrInvalidCredentials is returned when email or password checks fail.
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// AuthService handles high-level authentication processes.
type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*users.User, error)
	Login(ctx context.Context, req LoginRequest) (string, error)
}

type authService struct {
	userRepo users.UserRepository
	cfg      *config.Config
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo users.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// Register hashes the user password and creates a new database user.
func (s *authService) Register(ctx context.Context, req RegisterRequest) (*users.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &users.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login validates user credentials and issues a JWT token.
func (s *authService) Login(ctx context.Context, req LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	duration, err := time.ParseDuration(s.cfg.JWTExpiresIn)
	if err != nil {
		duration = 24 * time.Hour
	}

	token, err := GenerateToken(user.ID, s.cfg.JWTSecret, duration)
	if err != nil {
		return "", err
	}

	return token, nil
}
