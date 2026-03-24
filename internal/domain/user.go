package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID           uuid.UUID  `json:"id"               gorm:"primaryKey;type:char(36)"`
	Email        string     `json:"email"            gorm:"uniqueIndex;type:varchar(255);not null"`
	PasswordHash string     `json:"-"                gorm:"column:password_hash;type:varchar(255);not null"`
	Role         UserRole   `json:"role"             gorm:"type:enum('customer','admin');default:'customer';not null"`
	IsActive     bool       `json:"is_active"        gorm:"default:true;not null"`
	CreatedAt    time.Time  `json:"created_at"       gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at"       gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

func (User) TableName() string { return "users" }

// RegisterRequest holds input for user registration
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
	Phone     string `json:"phone" validate:"required"`
}

// LoginRequest holds input for login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse holds the JWT tokens returned after auth
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

// UserInfo is the safe public representation of a user
type UserInfo struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  UserRole  `json:"role"`
}

// RefreshTokenRequest holds the refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
