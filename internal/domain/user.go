package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RegisterRequest is the DTO for user registration
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// LoginRequest is the DTO for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse is the DTO returned after successful auth
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

const RoleAdmin = "admin"
const RoleUser = "user"

// UserRepository defines the interface for user data persistence
type (
	UserRepository interface {
		Create(ctx context.Context, user *User) error
		GetByID(ctx context.Context, id uuid.UUID) (*User, error)
		GetByEmail(ctx context.Context, email string) (*User, error)
		Update(ctx context.Context, user *User) error
		List(ctx context.Context) ([]User, error)
	}

	UserService interface {
		Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
		Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
		GetProfile(ctx context.Context, userID uuid.UUID) (*User, error)
	}

	AdminService interface {
		ListUsers(ctx context.Context) ([]User, error)
		ListAllTasks(ctx context.Context) ([]Task, error)
		ListUserTasks(ctx context.Context, userID uuid.UUID) ([]Task, error)
		PromoteToAdmin(ctx context.Context, userID uuid.UUID) (*User, error)
	}
)
