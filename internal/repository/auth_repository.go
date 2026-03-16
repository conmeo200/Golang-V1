package repository

import (
	"gorm.io/gorm"
)

// authRepository is the concrete implementation of AuthRepository.
type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new instance that implements AuthRepository.
func NewAuthRepository(db *gorm.DB) AuthRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	return &authRepository{db: db}
}
