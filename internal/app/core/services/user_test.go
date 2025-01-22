package services_test

import (
	"context"
	"github.com/bagashiz/go_hexagonal/internal/app/core/models"
	mock2 "github.com/bagashiz/go_hexagonal/internal/app/core/ports/mock"
	"github.com/bagashiz/go_hexagonal/internal/app/core/services"
	util2 "github.com/bagashiz/go_hexagonal/internal/app/core/utils"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type registerTestedInput struct {
	user *models.User
}

type registerExpectedOutput struct {
	user *models.User
	err  error
}

func TestUserService_Register(t *testing.T) {
	ctx := context.Background()
	userName := gofakeit.Name()
	userEmail := gofakeit.Email()
	userPassword := gofakeit.Password(true, true, true, true, false, 8)
	hashedPassword, _ := util2.HashPassword(userPassword)

	userInput := &models.User{
		Name:     userName,
		Email:    userEmail,
		Password: userPassword,
	}
	userOutput := &models.User{
		ID:        gofakeit.Uint64(),
		Name:      userName,
		Email:     userEmail,
		Password:  hashedPassword,
		Role:      models.Cashier,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cacheKey := util2.GenerateCacheKey("user", userOutput.ID)
	userSerialized, _ := util2.Serialize(userOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock2.MockUserRepository,
			cache *mock2.MockCacheRepository,
		)
		input    registerTestedInput
		expected registerExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(userOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("users:*")).
					Return(nil)
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: userOutput,
				err:  nil,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(nil, models.ErrInternal)
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: nil,
				err:  models.ErrInternal,
			},
		},
		{
			desc: "Fail_DuplicateData",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(nil, models.ErrConflictingData)
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: nil,
				err:  models.ErrConflictingData,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(userOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(models.ErrInternal)
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: nil,
				err:  models.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(userOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("users:*")).
					Return(models.ErrInternal)
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: nil,
				err:  models.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		// tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			// TODO: fix race condition to enable parallel testing
			// t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock2.NewMockUserRepository(ctrl)
			cache := mock2.NewMockCacheRepository(ctrl)

			tc.mocks(userRepo, cache)

			userService := services.NewUserService(userRepo, cache)

			user, err := userService.Register(ctx, tc.input.user)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.user, user, "User mismatch")
		})
	}
}

type getUserTestedInput struct {
	id uint64
}

type getUserExpectedOutput struct {
	user *models.User
	err  error
}

func TestUserService_GetUser(t *testing.T) {
	ctx := context.Background()
	userID := gofakeit.Uint64()
	userOutput := &models.User{
		ID:       userID,
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, 8),
		Role:     models.Cashier,
	}

	cacheKey := util2.GenerateCacheKey("user", userID)
	userSerialized, _ := util2.Serialize(userOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock2.MockUserRepository,
			cache *mock2.MockCacheRepository,
		)
		input    getUserTestedInput
		expected getUserExpectedOutput
	}{
		{
			desc: "Success_FromCache",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(userSerialized, nil)
			},
			input: getUserTestedInput{
				id: userID,
			},
			expected: getUserExpectedOutput{
				user: userOutput,
				err:  nil,
			},
		},
		{
			desc: "Success_FromDB",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, models.ErrDataNotFound)
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(userOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(nil)
			},
			input: getUserTestedInput{
				id: userID,
			},
			expected: getUserExpectedOutput{
				user: userOutput,
				err:  nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, models.ErrDataNotFound)
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(nil, models.ErrDataNotFound)
			},
			input: getUserTestedInput{
				id: userID,
			},
			expected: getUserExpectedOutput{
				user: nil,
				err:  models.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, models.ErrInternal)
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(nil, models.ErrInternal)
			},
			input: getUserTestedInput{
				id: userID,
			},
			expected: getUserExpectedOutput{
				user: nil,
				err:  models.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, models.ErrDataNotFound)
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(userOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(models.ErrInternal)
			},
			input: getUserTestedInput{
				id: userID,
			},
			expected: getUserExpectedOutput{
				user: nil,
				err:  models.ErrInternal,
			},
		},
		{
			desc: "Fail_Deserialize",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return([]byte("invalid"), nil)
			},
			input: getUserTestedInput{
				id: userID,
			},
			expected: getUserExpectedOutput{
				user: nil,
				err:  models.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock2.NewMockUserRepository(ctrl)
			cache := mock2.NewMockCacheRepository(ctrl)

			tc.mocks(userRepo, cache)

			userService := services.NewUserService(userRepo, cache)

			user, err := userService.GetUser(ctx, tc.input.id)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.user, user, "User mismatch")
		})
	}
}

type listUsersTestedInput struct {
	skip  uint64
	limit uint64
}

type listUsersExpectedOutput struct {
	users []models.User
	err   error
}

func TestUserService_ListUsers(t *testing.T) {
	var users []models.User

	for i := 0; i < 10; i++ {
		userPassword := gofakeit.Password(true, true, true, true, false, 8)
		hashedPassword, _ := util2.HashPassword(userPassword)

		users = append(users, models.User{
			ID:       gofakeit.Uint64(),
			Name:     gofakeit.Name(),
			Email:    gofakeit.Email(),
			Password: hashedPassword,
			Role:     models.Cashier,
		})
	}

	ctx := context.Background()
	skip := gofakeit.Uint64()
	limit := gofakeit.Uint64()

	params := util2.GenerateCacheKeyParams(skip, limit)
	cacheKey := util2.GenerateCacheKey("users", params)
	usersSerialized, _ := util2.Serialize(users)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock2.MockUserRepository,
			cache *mock2.MockCacheRepository,
		)
		input    listUsersTestedInput
		expected listUsersExpectedOutput
	}{
		{
			desc: "Success_FromCache",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(usersSerialized, nil)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: users,
				err:   nil,
			},
		},
		{
			desc: "Success_FromDB",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, models.ErrDataNotFound)
				userRepo.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(skip), gomock.Eq(limit)).
					Return(users, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(usersSerialized), gomock.Eq(ttl)).
					Return(nil)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: users,
				err:   nil,
			},
		},
		{
			desc: "Fail_Deserialize",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return([]byte("invalid"), nil)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: nil,
				err:   models.ErrInternal,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, models.ErrDataNotFound)
				userRepo.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(skip), gomock.Eq(limit)).
					Return(nil, models.ErrInternal)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: nil,
				err:   models.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, models.ErrDataNotFound)
				userRepo.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(skip), gomock.Eq(limit)).
					Return(users, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(usersSerialized), gomock.Eq(ttl)).
					Return(models.ErrInternal)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: nil,
				err:   models.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock2.NewMockUserRepository(ctrl)
			cache := mock2.NewMockCacheRepository(ctrl)

			tc.mocks(userRepo, cache)

			userService := services.NewUserService(userRepo, cache)

			users, err := userService.ListUsers(ctx, tc.input.skip, tc.input.limit)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.users, users, "Users mismatch")
		})
	}
}

type updateUserTestedInput struct {
	user *models.User
}

type updateUserExpectedOutput struct {
	user *models.User
	err  error
}

func TestUserService_UpdateUser(t *testing.T) {
	ctx := context.Background()
	userID := gofakeit.Uint64()

	// TODO: test with hashed password

	userInput := &models.User{
		ID:    userID,
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
		Role:  models.Cashier,
	}
	userOutput := &models.User{
		ID:    userID,
		Name:  userInput.Name,
		Email: userInput.Email,
		Role:  userInput.Role,
	}
	existingUser := &models.User{
		ID:    userID,
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
		Role:  models.Admin,
	}

	cacheKey := util2.GenerateCacheKey("user", userID)
	userSerialized, _ := util2.Serialize(userOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock2.MockUserRepository,
			cache *mock2.MockCacheRepository,
		)
		input    updateUserTestedInput
		expected updateUserExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
				userRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(userOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("users:*")).
					Return(nil)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: userOutput,
				err:  nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(nil, models.ErrDataNotFound)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  models.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalErrorGetByID",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(nil, models.ErrInternal)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  models.ErrInternal,
			},
		},
		{
			desc: "Fail_EmptyData",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
			},
			input: updateUserTestedInput{
				user: &models.User{
					ID: userID,
				},
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  models.ErrNoUpdatedData,
			},
		},
		{
			desc: "Fail_SameData",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
			},
			input: updateUserTestedInput{
				user: existingUser,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  models.ErrNoUpdatedData,
			},
		},
		{
			desc: "Fail_DuplicateData",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
				userRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(nil, models.ErrConflictingData)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  models.ErrConflictingData,
			},
		},
		{
			desc: "Fail_InternalErrorUpdate",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
				userRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(nil, models.ErrInternal)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  models.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCache",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
				userRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(userOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(models.ErrInternal)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  models.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
				userRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(userOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(models.ErrInternal)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  models.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
				userRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(userOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("users:*")).
					Return(models.ErrInternal)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  models.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		// tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			// TODO: fix race condition to enable parallel testing
			// t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock2.NewMockUserRepository(ctrl)
			cache := mock2.NewMockCacheRepository(ctrl)

			tc.mocks(userRepo, cache)

			userService := services.NewUserService(userRepo, cache)

			user, err := userService.UpdateUser(ctx, tc.input.user)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.user, user, "User mismatch")
		})
	}
}

type deleteUserTestedInput struct {
	id uint64
}

type deleteUserExpectedOutput struct {
	err error
}

func TestUserService_DeleteUser(t *testing.T) {
	ctx := context.Background()
	userID := gofakeit.Uint64()

	cacheKey := util2.GenerateCacheKey("user", userID)

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock2.MockUserRepository,
			cache *mock2.MockCacheRepository,
		)
		input    deleteUserTestedInput
		expected deleteUserExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(&models.User{}, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("users:*")).
					Return(nil)
				userRepo.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(userID)).
					Return(nil)
			},
			input: deleteUserTestedInput{
				id: userID,
			},
			expected: deleteUserExpectedOutput{
				err: nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(nil, models.ErrDataNotFound)
			},
			input: deleteUserTestedInput{
				id: userID,
			},
			expected: deleteUserExpectedOutput{
				err: models.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalErrorGetByID",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(nil, models.ErrInternal)
			},
			input: deleteUserTestedInput{
				id: userID,
			},
			expected: deleteUserExpectedOutput{
				err: models.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCache",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(&models.User{}, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(models.ErrInternal)
			},
			input: deleteUserTestedInput{
				id: userID,
			},
			expected: deleteUserExpectedOutput{
				err: models.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(&models.User{}, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("users:*")).
					Return(models.ErrInternal)
			},
			input: deleteUserTestedInput{
				id: userID,
			},
			expected: deleteUserExpectedOutput{
				err: models.ErrInternal,
			},
		},
		{
			desc: "Fail_InternalErrorDelete",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				cache *mock2.MockCacheRepository,
			) {
				user := &models.User{
					ID: userID,
				}
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(user, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("users:*")).
					Return(nil)
				userRepo.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(userID)).
					Return(models.ErrInternal)
			},
			input: deleteUserTestedInput{
				id: userID,
			},
			expected: deleteUserExpectedOutput{
				err: models.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock2.NewMockUserRepository(ctrl)
			cache := mock2.NewMockCacheRepository(ctrl)

			tc.mocks(userRepo, cache)

			userService := services.NewUserService(userRepo, cache)

			err := userService.DeleteUser(ctx, tc.input.id)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
		})
	}
}
