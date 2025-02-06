package models

import (
	"time"
)

// UserRole is an enum for user's role
type UserRole string

// User is an entity that represents a user
type User struct {
	ID        uint64
	Name      string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
