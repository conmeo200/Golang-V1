package service

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/dto"
)

// AuthService defines the interface for authentication-related business logic.
// It abstracts the implementation of login and registration processes.
type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (string, error)
	Register(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
}
