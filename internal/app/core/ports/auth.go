package ports

import (
	"context"
	"github.com/bagashiz/go_hexagonal/internal/app/core/models"
)

//go:generate mockgen -source=auth.go -destination=mock/auth.go -package=mock

// TokenService is an interface for interacting with token-related business logic
type TokenService interface {
	// CreateToken creates a new token for a given user
	CreateToken(user *models.User) (string, error)
	// VerifyToken verifies the token and returns the payload
	VerifyToken(token string) (*models.TokenPayload, error)
}

// UserService is an interface for interacting with user authentication-related business logic
type AuthService interface {
	// Login authenticates a user by email and password and returns a token
	Login(ctx context.Context, email, password string) (string, error)
}
