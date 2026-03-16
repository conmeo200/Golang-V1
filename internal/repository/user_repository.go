package repository

import (
	"errors"

	"github.com/conmeo200/Golang-V1/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {

	var user model.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *UserRepository) CreateUser(user *model.User) (*model.User, error) {

	result := r.db.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *UserRepository) GetUser(id string) (*model.User, error) {

	var user model.User

	err := r.db.First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *UserRepository) ListUser() ([]model.User, error) {

	var users []model.User

	err := r.db.Find(&users).Error

	return users, err
}

func (r *UserRepository) UpdateBalance(id uint, newBalance float64) error {

	return r.db.Model(&model.User{}).
		Where("id = ?", id).
		Update("balance", newBalance).Error
}

func (r *UserRepository) UpdatePassword(id string, newHash string) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", id).
		Update("password_hash", newHash).Error
}

func (r *UserRepository) Delete(id uint) error {

	return r.db.Delete(&model.User{}, id).Error
}