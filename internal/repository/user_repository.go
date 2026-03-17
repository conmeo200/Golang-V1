package repository

import (
	"context"
	"errors"

	"github.com/conmeo200/Golang-V1/internal/model"
	"gorm.io/gorm"
)

type UserRepo interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
	ListUser(ctx context.Context) ([]model.User, error)
	UpdateBalance(ctx context.Context, id uint, newBalance float64) error
	UpdatePassword(ctx context.Context, id string, newHash string) error
	Delete(ctx context.Context, id uint) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {

	var user model.User

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {

	result := r.db.WithContext(ctx).Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *UserRepository) GetUser(ctx context.Context, id string) (*model.User, error) {

	var user model.User

	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *UserRepository) ListUser(ctx context.Context) ([]model.User, error) {

	var users []model.User

	err := r.db.WithContext(ctx).Find(&users).Error

	return users, err
}

func (r *UserRepository) UpdateBalance(ctx context.Context, id uint, newBalance float64) error {

	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Update("balance", newBalance).Error
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id string, newHash string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Update("password_hash", newHash).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {

	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}