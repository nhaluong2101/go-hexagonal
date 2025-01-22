package auth

import (
	"github.com/bagashiz/go_hexagonal/internal/app/core/models"
	"github.com/bagashiz/go_hexagonal/internal/app/core/ports"
	"github.com/bagashiz/go_hexagonal/internal/app/infrastructure/configs"
	"go.uber.org/fx"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

/**
 * TokenHandler implements ports.TokenService interface
 * and provides an access to the paseto library
 */
type TokenHandler struct {
	token    *paseto.Token
	key      *paseto.V4SymmetricKey
	parser   *paseto.Parser
	duration time.Duration
}

// New creates a new paseto instance
func New(config *configs.Token) (ports.TokenService, error) {
	durationStr := config.Duration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, models.ErrTokenDuration
	}

	token := paseto.NewToken()
	key := paseto.NewV4SymmetricKey()
	parser := paseto.NewParser()

	return &TokenHandler{
		&token,
		&key,
		&parser,
		duration,
	}, nil
}

// CreateToken creates a new paseto token
func (pt *TokenHandler) CreateToken(user *models.User) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", models.ErrTokenCreation
	}

	payload := &models.TokenPayload{
		ID:     id,
		UserID: user.ID,
		Role:   user.Role,
	}

	err = pt.token.Set("payload", payload)
	if err != nil {
		return "", models.ErrTokenCreation
	}

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(pt.duration)

	pt.token.SetIssuedAt(issuedAt)
	pt.token.SetNotBefore(issuedAt)
	pt.token.SetExpiration(expiredAt)

	token := pt.token.V4Encrypt(*pt.key, nil)

	return token, nil
}

// VerifyToken verifies the paseto token
func (pt *TokenHandler) VerifyToken(token string) (*models.TokenPayload, error) {
	var payload *models.TokenPayload

	parsedToken, err := pt.parser.ParseV4Local(*pt.key, token, nil)
	if err != nil {
		if err.Error() == "this token has expired" {
			return nil, models.ErrExpiredToken
		}
		return nil, models.ErrInvalidToken
	}

	err = parsedToken.Get("payload", &payload)
	if err != nil {
		return nil, models.ErrInvalidToken
	}

	return payload, nil
}

var TokenModule = fx.Module(
	"token-handler-module",
	fx.Provide(
		fx.Annotate(New, fx.As(new(ports.TokenService))),
	),
)
