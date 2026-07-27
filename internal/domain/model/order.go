package model

import (
	"github.com/google/uuid"
)

type Order struct {
	UUID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index"`
	ShopID         uuid.UUID `gorm:"type:uuid;not null;index"`
	TotalAmount    float64   `gorm:"type:numeric(15,2);not null"`
	ShippingFee    float64   `gorm:"type:numeric(15,2);default:0"`
	DiscountAmount float64   `gorm:"type:numeric(15,2);default:0"`
	FinalAmount    float64   `gorm:"type:numeric(15,2);not null"`
	Status         string    `gorm:"type:varchar(50);default:'pending';not null;index;comment:'pending, processing, shipped, delivered, completed, cancelled'"`
	PaymentStatus  string    `gorm:"type:varchar(50);default:'unpaid';not null;index;comment:'unpaid, paid'"`
	PaymentMethod  string    `gorm:"type:varchar(50);default:'cod';not null"` // cod, vnpay, momo
	IdempotencyKey string    `gorm:"size:255;uniqueIndex"`
	ProcessedAt    int64
	CreatedAt      int64
	UpdatedAt      int64

	// Relationship
	User     User        `gorm:"foreignKey:UserID;references:ID"`
	Shop     Shop        `gorm:"foreignKey:ShopID"`
	Payment  []Payment   `gorm:"foreignKey:OrderID;references:UUID"`
	Items    []OrderItem `gorm:"foreignKey:OrderID;references:UUID"`
}

type OrderItem struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrderID          uuid.UUID `gorm:"type:uuid;index;not null"`
	ProductVariantID uuid.UUID `gorm:"type:uuid;index;not null"`
	Quantity         int       `gorm:"not null"`
	UnitPrice        float64   `gorm:"type:numeric(15,2);not null"`
	CreatedAt        int64
	UpdatedAt        int64
}
