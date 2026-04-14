package model

import (
	"github.com/google/uuid"
)

type Transaction struct {
	UUID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID        uuid.UUID `gorm:"type:uuid;not null;index"`
	PaymentID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Amount         float64   `gorm:"type:numeric(15,2);not null"`
	Status         string    `gorm:"type:varchar(50);default:'pending';not null;index;comment:'pending, success, failed'"`
	PaymentMethod  string    `gorm:"type:varchar(50);not null;comment:'vnpay, banking, cod'"`
	ReferenceID    string    `gorm:"size:255;index"`
	CreatedAt      int64
	UpdatedAt      int64

	// Relationship
	Order Order `gorm:"foreignKey:OrderID;references:UUID"`
	Payment Payment `gorm:"foreignKey:PaymentID;references:UUID"`
}
