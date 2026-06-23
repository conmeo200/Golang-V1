package model

import (
	"github.com/google/uuid"
)

type WebhookLog struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Provider string    `gorm:"type:varchar(50);index"`
	EventID  string    `gorm:"type:varchar(255);index"`

	Payload []byte `gorm:"type:jsonb"`

	Signature string `gorm:"type:text"`
	Verified  bool   `gorm:"default:false"`

	Status       string `gorm:"type:varchar(20);index"`
	ErrorMessage string `gorm:"type:text"`

	CreatedAt int64 `gorm:"index"`
}