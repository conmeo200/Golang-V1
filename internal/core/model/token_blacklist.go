package model

import "github.com/google/uuid"

type TokenBlacklist struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Token     string    `gorm:"type:text;uniqueIndex;not null"`
	ExpiresAt int64     `gorm:"not null"`
	CreatedAt int64     `gorm:"autoCreateTime"`
}
