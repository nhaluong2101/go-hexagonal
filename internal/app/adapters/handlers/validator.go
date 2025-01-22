package handlers

import (
	"github.com/bagashiz/go_hexagonal/internal/app/core/models"
	"github.com/go-playground/validator/v10"
)

// userRoleValidator is a custom validator for validating user roles
var userRoleValidator validator.Func = func(fl validator.FieldLevel) bool {
	userRole := fl.Field().Interface().(models.UserRole)

	switch userRole {
	case "admin", "cashier":
		return true
	default:
		return false
	}
}
