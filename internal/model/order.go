package model

import (
	"github.com/google/uuid"
)

type Order struct {
	UUID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index"`
	Amount         float64   `gorm:"type:numeric(15,2);not null"`
	Status         string    `gorm:"type:varchar(50);default:'pending';not null;index;comment:'pending, processing, completed, failed'"`
	PaymentStatus  string    `gorm:"type:varchar(50);default:'unpaid';not null;index;comment:'unpaid, paid'"`
	IdempotencyKey string    `gorm:"size:255;uniqueIndex"`
	ProcessedAt    int64
	CreatedAt      int64
	UpdatedAt      int64

	// Relationship
	User     User      `gorm:"foreignKey:UserID;references:ID"`
	Payment []Payment `gorm:"foreignKey:OrderID;references:UUID"`
}
