package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	InternalID int64          `json:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID      `json:"public_id"   gorm:"column:public_id;type:uuid"`
	Name       string         `json:"name"        gorm:"column:name"`
	Email      string         `json:"email"       gorm:"column:email;unique"`
	Password   string         `json:"-"           gorm:"column:password"`
	Role       string         `json:"role"        gorm:"column:role"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-"           gorm:"index"`
}

// UserResponse is the public-facing user object (no password).
type UserResponse struct {
	PublicID  uuid.UUID `json:"public_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterRequest is the body for POST /v1/auth/register.
type RegisterRequest struct {
	Name     string `json:"name"     example:"John Doe"`
	Email    string `json:"email"    example:"john@example.com"`
	Password string `json:"password" example:"secret123"`
}

// LoginRequest is the body for POST /v1/auth/login.
type LoginRequest struct {
	Email    string `json:"email"    example:"admin@gmail.com"`
	Password string `json:"password" example:"admin123"`
}

// LoginResponse is returned on successful login.
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// UpdateUserRequest is the body for PUT /api/v1/users/:id.
type UpdateUserRequest struct {
	Name string `json:"name" example:"John Updated"`
}
