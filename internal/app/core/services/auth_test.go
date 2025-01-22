package services_test

import (
	"context"
	"github.com/bagashiz/go_hexagonal/internal/app/core/models"
	mock2 "github.com/bagashiz/go_hexagonal/internal/app/core/ports/mock"
	"github.com/bagashiz/go_hexagonal/internal/app/core/services"
	"github.com/bagashiz/go_hexagonal/internal/app/core/utils"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"go.uber.org/mock/gomock"
)

type loginTestedInput struct {
	email    string
	password string
}

type loginExpectedOutput struct {
	token string
	err   error
}

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 8)
	hashedPassword, _ := utils.HashPassword(password)
	user := &models.User{
		Email:    email,
		Password: hashedPassword,
	}
	failUser := &models.User{
		Email:    email,
		Password: "wrong password",
	}
	token := gofakeit.UUID()

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock2.MockUserRepository,
			tokenService *mock2.MockTokenService,
		)
		input    loginTestedInput
		expected loginExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				tokenService *mock2.MockTokenService,
			) {
				userRepo.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(email)).
					Times(1).
					Return(user, nil)
				tokenService.EXPECT().
					CreateToken(gomock.Eq(user)).
					Times(1).
					Return(token, nil)
			},
			input: loginTestedInput{
				email:    email,
				password: password,
			},
			expected: loginExpectedOutput{
				token: token,
				err:   nil,
			},
		},
		{
			desc: "Fail_UserNotFound",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				tokenService *mock2.MockTokenService,
			) {
				userRepo.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(email)).
					Times(1).
					Return(nil, models.ErrDataNotFound)
			},
			input: loginTestedInput{
				email:    email,
				password: password,
			},
			expected: loginExpectedOutput{
				token: "",
				err:   models.ErrInvalidCredentials,
			},
		},
		{
			desc: "Fail_PasswordMismatch",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				tokenService *mock2.MockTokenService,
			) {
				userRepo.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(email)).
					Times(1).
					Return(failUser, nil)
			},
			input: loginTestedInput{
				email:    email,
				password: password,
			},
			expected: loginExpectedOutput{
				token: "",
				err:   models.ErrInvalidCredentials,
			},
		},
		{
			desc: "Fail_TokenCreation",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				tokenService *mock2.MockTokenService,
			) {
				userRepo.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(email)).
					Times(1).
					Return(user, nil)
				tokenService.EXPECT().
					CreateToken(gomock.Eq(user)).
					Times(1).
					Return("", models.ErrTokenCreation)
			},
			input: loginTestedInput{
				email:    email,
				password: password,
			},
			expected: loginExpectedOutput{
				token: "",
				err:   models.ErrTokenCreation,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				userRepo *mock2.MockUserRepository,
				tokenService *mock2.MockTokenService,
			) {
				userRepo.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(email)).
					Times(1).
					Return(nil, models.ErrInternal)
			},
			input: loginTestedInput{
				email:    email,
				password: password,
			},
			expected: loginExpectedOutput{
				token: "",
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
			tokenService := mock2.NewMockTokenService(ctrl)

			tc.mocks(userRepo, tokenService)

			authService := services.NewAuthService(userRepo, tokenService)

			token, err := authService.Login(ctx, tc.input.email, tc.input.password)
			if err != tc.expected.err {
				t.Errorf("[case: %s] expected to get %q; got %q", tc.desc, tc.expected.err, err)
			}
			if token != tc.expected.token {
				t.Errorf("[case: %s] expected to get %q; got %q", tc.desc, tc.expected.token, token)
			}
		})
	}
}
