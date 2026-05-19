package model

import (
	"github.com/google/uuid"
)

type DeadLetterEvent struct {
	ID           uint      `gorm:"primaryKey"`
	EventID      uuid.UUID `gorm:"type:uuid;index"`
	QueueName    string    `gorm:"type:varchar(255);index"`
	ExchangeName string    `gorm:"type:varchar(255)"`
	RoutingKey   string    `gorm:"type:varchar(255)"`
	Payload      []byte    `gorm:"type:text"`
	Headers      []byte    `gorm:"type:text"` // JSON encoded headers
	Reason       string    `gorm:"type:text"`
	Status       string    `gorm:"type:varchar(50);default:'pending';index"` // pending, replayed, resolved
	RetryCount   int       `gorm:"default:0"`
	CreatedAt    int64     `gorm:"autoCreateTime"`
	UpdatedAt    int64     `gorm:"autoUpdateTime"`
}

func (DeadLetterEvent) TableName() string {
	return "dead_letter_events"
}

type DeadLetterRepo interface {
	Create(event *DeadLetterEvent) error
	GetByID(id uint) (*DeadLetterEvent, error)
	ListPending() ([]DeadLetterEvent, error)
	MarkAsReplayed(id uint) error
	MarkAsResolved(id uint) error
	Exists(eventID uuid.UUID, queueName string) (bool, error)
}
