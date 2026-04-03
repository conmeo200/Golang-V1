package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	Role         string  `gorm:"size:50;default:user;index"`
	RoleID       uint    `gorm:"index" json:"role_id"`
	RoleEntity   Role    `gorm:"foreignKey:RoleID" json:"role_entity"`
	Status       string  `gorm:"size:50;default:active;index"`
	Balance      float64 `gorm:"type:numeric(15,2);default:0"`
	LastLoginAt int64
	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
}