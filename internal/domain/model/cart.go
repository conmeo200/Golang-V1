package model

import (
	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	CreatedAt int64
	UpdatedAt int64
	
	Items     []CartItem `gorm:"foreignKey:CartID"`
}

type CartItem struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CartID           uuid.UUID `gorm:"type:uuid;index;not null"`
	ProductVariantID uuid.UUID `gorm:"type:uuid;index;not null"`
	Quantity         int       `gorm:"not null;default:1"`
	CreatedAt        int64
	UpdatedAt        int64
}

type ShippingAddress struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      uuid.UUID `gorm:"type:uuid;index;not null"`
	FullName    string    `gorm:"size:255;not null"`
	Phone       string    `gorm:"size:20;not null"`
	AddressLine string    `gorm:"size:255;not null"`
	City        string    `gorm:"size:100;not null"`
	District    string    `gorm:"size:100;not null"`
	Ward        string    `gorm:"size:100;not null"`
	IsDefault   bool      `gorm:"default:false"`
	CreatedAt   int64
	UpdatedAt   int64
}
