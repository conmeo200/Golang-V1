package seeder

import (
	"time"

	"github.com/conmeo200/Golang-V1/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {

	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	now := time.Now().Unix()

	users := []model.User{
		{
			Email:        "admin@example.com",
			PasswordHash: string(hash),
			Role:         "admin",
			Status:       "active",
			Balance:      1000,
			LastLoginAt:  now,
			CreatedAt:    now,
			UpdatedAt:    0,
		},
		{
			Email:        "user@example.com",
			PasswordHash: string(hash),
			Role:         "user",
			Status:       "active",
			Balance:      100,
			LastLoginAt:   now,
			CreatedAt:    now,
			UpdatedAt:    0,
		},
	}

	return db.Create(&users).Error
}