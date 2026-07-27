package model

import (
	"github.com/google/uuid"
)

type Shop struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OwnerID     uuid.UUID `gorm:"type:uuid;index;not null"`
	Owner       User      `gorm:"foreignKey:OwnerID"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	LogoURL     string    `gorm:"size:255"`
	Status      string    `gorm:"size:50;default:active;index"` // active, banned, suspended
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   int64     `gorm:"index"`
}

type ShopAddress struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ShopID      uuid.UUID `gorm:"type:uuid;index;not null"`
	AddressLine string    `gorm:"size:255;not null"`
	City        string    `gorm:"size:100;not null"`
	District    string    `gorm:"size:100;not null"`
	Ward        string    `gorm:"size:100;not null"`
	IsDefault   bool      `gorm:"default:false"`
	CreatedAt   int64
	UpdatedAt   int64
}
