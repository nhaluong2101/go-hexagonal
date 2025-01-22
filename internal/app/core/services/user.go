package services

import (
	"context"
	"github.com/bagashiz/go_hexagonal/internal/app/core/models"
	"github.com/bagashiz/go_hexagonal/internal/app/core/ports"
	"github.com/bagashiz/go_hexagonal/internal/app/core/utils"
)

/**
 * UserService implements ports.UserService interface
 * and provides an access to the user repositories
 * and cache services
 */
type UserService struct {
	repo  ports.UserRepository
	cache ports.CacheRepository
}

// NewUserService creates a new user services instance
func NewUserService(repo ports.UserRepository, cache ports.CacheRepository) *UserService {
	return &UserService{
		repo,
		cache,
	}
}

// Register creates a new user
func (us *UserService) Register(ctx context.Context, user *models.User) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, models.ErrInternal
	}

	user.Password = hashedPassword

	user, err = us.repo.CreateUser(ctx, user)
	if err != nil {
		if err == models.ErrConflictingData {
			return nil, err
		}
		return nil, models.ErrInternal
	}

	cacheKey := utils.GenerateCacheKey("user", user.ID)
	userSerialized, err := utils.Serialize(user)
	if err != nil {
		return nil, models.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, models.ErrInternal
	}

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, models.ErrInternal
	}

	return user, nil
}

// GetUser gets a user by ID
func (us *UserService) GetUser(ctx context.Context, id uint64) (*models.User, error) {
	var user *models.User

	cacheKey := utils.GenerateCacheKey("user", id)
	cachedUser, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		err := utils.Deserialize(cachedUser, &user)
		if err != nil {
			return nil, models.ErrInternal
		}
		return user, nil
	}

	user, err = us.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == models.ErrDataNotFound {
			return nil, err
		}
		return nil, models.ErrInternal
	}

	userSerialized, err := utils.Serialize(user)
	if err != nil {
		return nil, models.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, models.ErrInternal
	}

	return user, nil
}

// ListUsers lists all users
func (us *UserService) ListUsers(ctx context.Context, skip, limit uint64) ([]models.User, error) {
	var users []models.User

	params := utils.GenerateCacheKeyParams(skip, limit)
	cacheKey := utils.GenerateCacheKey("users", params)

	cachedUsers, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		err := utils.Deserialize(cachedUsers, &users)
		if err != nil {
			return nil, models.ErrInternal
		}
		return users, nil
	}

	users, err = us.repo.ListUsers(ctx, skip, limit)
	if err != nil {
		return nil, models.ErrInternal
	}

	usersSerialized, err := utils.Serialize(users)
	if err != nil {
		return nil, models.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, usersSerialized, 0)
	if err != nil {
		return nil, models.ErrInternal
	}

	return users, nil
}

// UpdateUser updates a user's name, email, and password
func (us *UserService) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	existingUser, err := us.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		if err == models.ErrDataNotFound {
			return nil, err
		}
		return nil, models.ErrInternal
	}

	emptyData := user.Name == "" &&
		user.Email == "" &&
		user.Password == "" &&
		user.Role == ""
	sameData := existingUser.Name == user.Name &&
		existingUser.Email == user.Email &&
		existingUser.Role == user.Role
	if emptyData || sameData {
		return nil, models.ErrNoUpdatedData
	}

	var hashedPassword string

	if user.Password != "" {
		hashedPassword, err = utils.HashPassword(user.Password)
		if err != nil {
			return nil, models.ErrInternal
		}
	}

	user.Password = hashedPassword

	_, err = us.repo.UpdateUser(ctx, user)
	if err != nil {
		if err == models.ErrConflictingData {
			return nil, err
		}
		return nil, models.ErrInternal
	}

	cacheKey := utils.GenerateCacheKey("user", user.ID)

	err = us.cache.Delete(ctx, cacheKey)
	if err != nil {
		return nil, models.ErrInternal
	}

	userSerialized, err := utils.Serialize(user)
	if err != nil {
		return nil, models.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, models.ErrInternal
	}

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, models.ErrInternal
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (us *UserService) DeleteUser(ctx context.Context, id uint64) error {
	_, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == models.ErrDataNotFound {
			return err
		}
		return models.ErrInternal
	}

	cacheKey := utils.GenerateCacheKey("user", id)

	err = us.cache.Delete(ctx, cacheKey)
	if err != nil {
		return models.ErrInternal
	}

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return models.ErrInternal
	}

	return us.repo.DeleteUser(ctx, id)
}
