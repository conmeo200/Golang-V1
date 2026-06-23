package port

import (
	"context"
	"github.com/conmeo200/Golang-V1/internal/domain/model"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
	ListUser(ctx context.Context) ([]model.User, error)
	UpdateBalance(ctx context.Context, id uint, newBalance float64) error
	UpdatePassword(ctx context.Context, id string, newHash string) error
	Delete(ctx context.Context, id uint) error
}
