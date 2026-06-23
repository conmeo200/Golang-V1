package port

import (
	"context"
	"github.com/conmeo200/Golang-V1/internal/domain/model"
)

type UserService interface {
	FindFirstByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, email string, balance float64, password string) (*model.User, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
	ListUser(ctx context.Context) ([]model.User, error)
	UpdateBalance(ctx context.Context, id uint, newBalance float64) error
	DeleteUser(ctx context.Context, id uint) error
}
