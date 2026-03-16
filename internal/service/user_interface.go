package service

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/dto"
	"github.com/conmeo200/Golang-V1/internal/model"
)

// UserService defines the interface for user-related business logic.
// It abstracts the underlying implementation of the user service,
// allowing for easier testing and dependency injection.
type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
}
