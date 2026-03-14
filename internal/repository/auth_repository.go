package repository

import (
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	return &AuthRepository{db: db}
}

