package ports

import (
	"context"
	"github.com/bagashiz/go_hexagonal/internal/app/core/models"
)

//go:generate mockgen -source=user.go -destination=mock/user.go -package=mock

// UserRepository is an interface for interacting with user-related data
type UserRepository interface {
	// CreateUser inserts a new user into the database
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	// GetUserByID selects a user by id
	GetUserByID(ctx context.Context, id uint64) (*models.User, error)
	// GetUserByEmail selects a user by email
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	// ListUsers selects a list of users with pagination
	ListUsers(ctx context.Context, skip, limit uint64) ([]models.User, error)
	// UpdateUser updates a user
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id uint64) error
}

// UserService is an interface for interacting with user-related business logic
type UserService interface {
	// Register registers a new user
	Register(ctx context.Context, user *models.User) (*models.User, error)
	// GetUser returns a user by id
	GetUser(ctx context.Context, id uint64) (*models.User, error)
	// ListUsers returns a list of users with pagination
	ListUsers(ctx context.Context, skip, limit uint64) ([]models.User, error)
	// UpdateUser updates a user
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id uint64) error
}
