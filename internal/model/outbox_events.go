package model

import (
	"github.com/google/uuid"
)

type OutboxEvents struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	EventID   uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	EventType string    `gorm:"type:varchar(100);not null;index"`
	Payload   []byte    `gorm:"type:jsonb"`

	Status     string `gorm:"type:varchar(20);default:'PENDING';index;comment:'PENDING, SENT, FAILED'"`
	RetryCount int    `gorm:"default:0"`

	CreatedAt   int64 `gorm:"index"`
	SentAt      *int64
	NextRetryAt int64 `gorm:"index"`
}
