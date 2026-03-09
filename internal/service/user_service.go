package service

import (
	"github.com/conmeo200/Golang-V1/internal/model"

	"gorm.io/gorm"
)
type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	if db == nil {
        panic("db cannot be nil")
    }
	
	return &UserService{db: db}
}

func (s *UserService) CreateUser(email string, balance float64) error {
	user := model.User{
		Email:   email,
		Balance: balance,
	}
	return s.db.Create(&user).Error
}

func (s *UserService) GetUser(id uint) (*model.User, error) {
	var user model.User
	err := s.db.First(&user, id).Error
	return &user, err
}
func (s *UserService) ListUser() ([]*model.User, error) {
	var users []*model.User
	err := s.db.Find(&users).Error
	return users, err
}

func (s *UserService) Update(id uint, newBalance float64) error {
	return s.db.Model(&model.User{}).
		Where("id = ?", id).
		Update("balance", newBalance).Error
}

func (s *UserService) Delete(id uint) error {
	return s.db.Delete(&model.User{}, id).Error
}