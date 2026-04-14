package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"

)


type PaymentEvent struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PaymentID uuid.UUID `gorm:"type:uuid;not null;index"`

	EventType string `gorm:"type:varchar(50);not null;index"`
	Payload   datatypes.JSON `gorm:"type:jsonb"`


	Source string `gorm:"type:varchar(50);comment:'system, webhook'"`

	CreatedAt int64

	// Relationships
	Payment Payment `gorm:"foreignKey:PaymentID;references:UUID"`
}