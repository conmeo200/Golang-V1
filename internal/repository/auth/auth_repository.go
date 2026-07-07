package repository

import (
	"github.com/conmeo200/Golang-V1/internal/domain/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository interface {
	ProcessLogin() (model.User, error)
	ProcessRegister() (model.User, error)
	ProcessLogout() (uuid.UUID, error)
	// ProcessRefresh() (model.User, error)
	// ProcessUpdatePassword() (model.User, error)
	// ProcessUpdateBalance() (model.User, error)
	// ProcessDeleteUser() (uuid.UUID, error)
}

type AuthRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{db: db}
}

func (r *AuthRepositoryImpl) ProcessLogin() (model.User, error) {
	return model.User{}, nil
}

func (r *AuthRepositoryImpl) ProcessRegister() (model.User, error) {
	return model.User{}, nil
}

func (r *AuthRepositoryImpl) ProcessLogout() (uuid.UUID, error) {
	return uuid.Nil, nil
}

