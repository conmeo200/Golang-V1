package model

import (
	"time"
)

type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:100;not null;uniqueIndex"`
	Description string `gorm:"size:255"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`

	// Relationships
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}
