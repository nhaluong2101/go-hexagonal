package utils

import (
	"errors"
	"github.com/bagashiz/go_hexagonal/internal/app/core/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

// response represents a response body format
type response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

// NewResponse is a helper function to create a response body
func NewResponse(success bool, message string, data any) response {
	return response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

// meta represents metadata for a paginated response
type Meta struct {
	Total uint64 `json:"total" example:"100"`
	Limit uint64 `json:"limit" example:"10"`
	Skip  uint64 `json:"skip" example:"0"`
}

// NewMeta is a helper function to create metadata for a paginated response
func NewMeta(total, limit, skip uint64) Meta {
	return Meta{
		Total: total,
		Limit: limit,
		Skip:  skip,
	}
}

// authResponse represents an authentication response body
type authResponse struct {
	AccessToken string `json:"token" example:"v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..."`
}

// NewAuthResponse is a helper function to create a response body for handling authentication data
func NewAuthResponse(token string) authResponse {
	return authResponse{
		AccessToken: token,
	}
}

// userResponse represents a user response body
type UserResponse struct {
	ID        uint64    `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"test@example.com"`
	CreatedAt time.Time `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// NewUserResponse is a helper function to create a response body for handling user data
func NewUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// errorStatusMap is a map of defined error messages and their corresponding http status codes
var errorStatusMap = map[error]int{
	models.ErrInternal:                   http.StatusInternalServerError,
	models.ErrDataNotFound:               http.StatusNotFound,
	models.ErrConflictingData:            http.StatusConflict,
	models.ErrInvalidCredentials:         http.StatusUnauthorized,
	models.ErrUnauthorized:               http.StatusUnauthorized,
	models.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	models.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	models.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	models.ErrInvalidToken:               http.StatusUnauthorized,
	models.ErrExpiredToken:               http.StatusUnauthorized,
	models.ErrForbidden:                  http.StatusForbidden,
	models.ErrNoUpdatedData:              http.StatusBadRequest,
	models.ErrInsufficientStock:          http.StatusBadRequest,
	models.ErrInsufficientPayment:        http.StatusBadRequest,
}

// ValidationError sends an error response for some specific request validation error
func ValidationError(ctx *gin.Context, err error) {
	errMsgs := ParseError(err)
	errRsp := NewErrorResponse(errMsgs)
	ctx.JSON(http.StatusBadRequest, errRsp)
}

// HandleError determines the status code of an error and returns a JSON response with the error message and status code
func HandleError(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := ParseError(err)
	errRsp := NewErrorResponse(errMsg)
	ctx.JSON(statusCode, errRsp)
}

// HandleAbort sends an error response and aborts the request with the specified status code and error message
func HandleAbort(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := ParseError(err)
	errRsp := NewErrorResponse(errMsg)
	ctx.AbortWithStatusJSON(statusCode, errRsp)
}

// ParseError parses error messages from the error object and returns a slice of error messages
func ParseError(err error) []string {
	var errMsgs []string

	if errors.As(err, &validator.ValidationErrors{}) {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, err.Error())
		}
	} else {
		errMsgs = append(errMsgs, err.Error())
	}

	return errMsgs
}

// errorResponse represents an error response body format
type errorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

// NewErrorResponse is a helper function to create an error response body
func NewErrorResponse(errMsgs []string) errorResponse {
	return errorResponse{
		Success:  false,
		Messages: errMsgs,
	}
}

// HandleSuccess sends a success response with the specified status code and optional data
func HandleSuccess(ctx *gin.Context, data any) {
	rsp := NewResponse(true, "Success", data)
	ctx.JSON(http.StatusOK, rsp)
}
