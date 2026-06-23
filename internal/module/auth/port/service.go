package port

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/domain/model"
)

type AuthService interface {
	RegisterUser(ctx context.Context, email string, password string) (*model.User, error)
	LoginUser(ctx context.Context, email string, password string) (*model.User, error)
	RevokeToken(ctx context.Context, tokenString string, expiresAt int64) error
	IsTokenBlacklisted(ctx context.Context, tokenString string) bool
	ChangePassword(ctx context.Context, userID string, oldPassword string, newPassword string) error
	ForgotPassword(ctx context.Context, email string) (string, error)
	RefreshToken(ctx context.Context, tokenString string) (string, string, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}
