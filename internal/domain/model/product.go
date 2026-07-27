package model

import (
	"github.com/google/uuid"
)

type ProductCategory struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ParentID  *uuid.UUID `gorm:"type:uuid;index"`
	Name      string    `gorm:"size:255;not null"`
	Slug      string    `gorm:"size:255;uniqueIndex;not null"`
	IconURL   string    `gorm:"size:255"`
	CreatedAt int64
	UpdatedAt int64
}

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ShopID      uuid.UUID `gorm:"type:uuid;index;not null"`
	CategoryID  uuid.UUID `gorm:"type:uuid;index;not null"`
	Name        string    `gorm:"size:255;not null"`
	Slug        string    `gorm:"size:255;uniqueIndex;not null"`
	Description string    `gorm:"type:text"`
	BasePrice   float64   `gorm:"type:numeric(15,2);not null"`
	IsActive    bool      `gorm:"default:true;index"`
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   int64     `gorm:"index"`

	// Associations
	Variants []ProductVariant `gorm:"foreignKey:ProductID"`
	Images   []ProductImage   `gorm:"foreignKey:ProductID"`
}

type ProductVariant struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProductID uuid.UUID `gorm:"type:uuid;index;not null"`
	SKU       string    `gorm:"size:100;uniqueIndex;not null"`
	Name      string    `gorm:"size:255;not null"` // e.g. "Color: Red, Size: L"
	Price     float64   `gorm:"type:numeric(15,2);not null"`
	CreatedAt int64
	UpdatedAt int64
}

type ProductImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProductID uuid.UUID `gorm:"type:uuid;index;not null"`
	ImageURL  string    `gorm:"size:255;not null"`
	IsPrimary bool      `gorm:"default:false"`
	CreatedAt int64
}
