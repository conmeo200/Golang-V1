package database

import (
	"github.com/conmeo200/Golang-V1/internal/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.TokenBlacklist{},
		&model.Order{},
		&model.Role{},
		&model.Permission{},
	)
}
