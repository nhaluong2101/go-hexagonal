package services

import (
	"context"
	"github.com/bagashiz/go_hexagonal/internal/app/core/models"
	"github.com/bagashiz/go_hexagonal/internal/app/core/ports"
	"github.com/bagashiz/go_hexagonal/internal/app/core/utils"
)

/**
 * AuthService implements ports.AuthService interface
 * and provides an access to the user repositories
 * and token services
 */
type AuthService struct {
	repo ports.UserRepository
	ts   ports.TokenService
}

// NewAuthService creates a new auth services instance
func NewAuthService(repo ports.UserRepository, ts ports.TokenService) *AuthService {
	return &AuthService{
		repo,
		ts,
	}
}

// Login gives a registered user an access token if the credentials are valid
func (as *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := as.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == models.ErrDataNotFound {
			return "", models.ErrInvalidCredentials
		}
		return "", models.ErrInternal
	}

	err = utils.ComparePassword(password, user.Password)
	if err != nil {
		return "", models.ErrInvalidCredentials
	}

	accessToken, err := as.ts.CreateToken(user)
	if err != nil {
		return "", models.ErrTokenCreation
	}

	return accessToken, nil
}
