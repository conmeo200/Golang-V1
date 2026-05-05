package model

import (
	"github.com/google/uuid"
)

type Payment struct {
	UUID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID uuid.UUID `gorm:"type:uuid;not null;index"`

	Amount   float64 `gorm:"type:numeric(15,2);not null"`
	Currency string  `gorm:"type:varchar(10);not null"`

	Provider          string `gorm:"type:varchar(50);index"`
	ProviderPaymentID string `gorm:"type:varchar(255);index"`

	Status        string `gorm:"type:varchar(50);index;comment:'pending, authorized, captured, failed, refunded'"`
	PaymentMethod string `gorm:"type:varchar(50)"`

	IdempotencyKey string `gorm:"type:varchar(255);uniqueIndex"`

	CreatedAt int64
	UpdatedAt int64

	// Relationships
	Order  Order          `gorm:"foreignKey:OrderID;references:UUID"`
	Events []PaymentEvent `gorm:"foreignKey:PaymentID"`
}