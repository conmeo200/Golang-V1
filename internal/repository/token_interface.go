package repository

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/model"
)

// TokenRepository defines the interface for operations related to token management,
// such as blacklisting and checking the blacklist.
type TokenRepository interface {
	BlacklistToken(ctx context.Context, token *model.TokenBlacklist) error
	IsBlacklisted(ctx context.Context, tokenString string) bool
}
