package model

import (
	"github.com/google/uuid"
)

type InboxEvent struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	EventID   uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	EventType string    `gorm:"type:varchar(100);not null"`
	Payload   []byte    `gorm:"type:jsonb"`

	Status string `gorm:"type:varchar(20);default:'PROCESSED';index;comment:'PROCESSED, FAILED'"`

	CreatedAt   int64 `gorm:"index"`
	ProcessedAt int64
}
