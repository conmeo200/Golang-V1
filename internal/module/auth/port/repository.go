package port

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/domain/model"
)

type TokenRepository interface {
	BlacklistToken(ctx context.Context, token *model.TokenBlacklist) error
	IsBlacklisted(ctx context.Context, tokenString string) bool
}

// AuthRepository can be expanded later
type AuthRepository interface {
}
