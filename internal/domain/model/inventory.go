package model

import (
	"github.com/google/uuid"
)

type Inventory struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProductVariantID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	AvailableStock   int       `gorm:"default:0;not null"`
	ReservedStock    int       `gorm:"default:0;not null"` // Stock reserved for unpaid orders
	CreatedAt        int64
	UpdatedAt        int64
}

type InventoryLog struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	InventoryID      uuid.UUID `gorm:"type:uuid;index;not null"`
	OrderID          *uuid.UUID `gorm:"type:uuid;index"` // Nullable if the stock change is not due to an order (e.g. manual restock)
	ChangeAmount     int       `gorm:"not null"` // Positive for restock, negative for sale
	Reason           string    `gorm:"size:255;not null"` // e.g. "Order Placed", "Restock", "Adjustment"
	CreatedAt        int64
}
