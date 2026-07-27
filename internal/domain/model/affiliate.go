package model

import (
	"github.com/google/uuid"
)

type AffiliateProfile struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID          uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	AffiliateCode   string    `gorm:"size:50;uniqueIndex;not null"`
	TotalCommission float64   `gorm:"type:numeric(15,2);default:0"`
	Balance         float64   `gorm:"type:numeric(15,2);default:0"`
	CreatedAt       int64
	UpdatedAt       int64
	
	Links       []AffiliateLink       `gorm:"foreignKey:AffiliateProfileID"`
	Commissions []AffiliateCommission `gorm:"foreignKey:AffiliateProfileID"`
}

type AffiliateLink struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	AffiliateProfileID uuid.UUID `gorm:"type:uuid;index;not null"`
	OriginalURL        string    `gorm:"size:1024;not null"`
	ShortCode          string    `gorm:"size:20;uniqueIndex;not null"`
	Clicks             int       `gorm:"default:0"`
	CreatedAt          int64
}

type AffiliateCommission struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	AffiliateProfileID uuid.UUID `gorm:"type:uuid;index;not null"`
	OrderID            uuid.UUID `gorm:"type:uuid;index;not null"`
	Amount             float64   `gorm:"type:numeric(15,2);not null"`
	Status             string    `gorm:"size:50;default:pending;index"` // pending, approved, paid, cancelled
	CreatedAt          int64
	UpdatedAt          int64
}
