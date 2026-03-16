package repository

import (
	"context"
	"errors"

	"github.com/conmeo200/Golang-V1/internal/model"
	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	return &TokenRepository{db: db}
}

func (r *TokenRepository) BlacklistToken(ctx context.Context, token *model.TokenBlacklist) error {
	result := r.db.WithContext(ctx).Create(token)
	return result.Error
}

func (r *TokenRepository) IsBlacklisted(ctx context.Context, tokenString string) bool {
	var token model.TokenBlacklist
	err := r.db.WithContext(ctx).Where("token = ?", tokenString).First(&token).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
