package handlers

import (
	"github.com/bagashiz/go_hexagonal/internal/app/core/ports"
	"github.com/bagashiz/go_hexagonal/internal/app/core/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// AuthHandler represents the HTTP handlers for authentication-related requests
type AuthHandler struct {
	svc ports.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(svc ports.AuthService) *AuthHandler {
	return &AuthHandler{
		svc,
	}
}

// loginRequest represents the request body for logging in a user
type loginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678" minLength:"8"`
}

// Login godoc
//
//	@Summary		Login and get an access token
//	@Description	Logs in a registered user and returns an access token if the credentials are valid.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loginRequest	true	"Login request body"
//	@Success		200		{object}	authResponse	"Succesfully logged in"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		401		{object}	errorResponse	"Unauthorized error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/users/login [post]
func (ah *AuthHandler) Login(ctx *gin.Context) {
	var req loginRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	token, err := ah.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	rsp := utils.NewAuthResponse(token)

	utils.HandleSuccess(ctx, rsp)
}

var AuthModule = fx.Module(
	"auth-handler-module",
	fx.Provide(NewAuthHandler),
)
