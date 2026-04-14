package repository

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/model"
	"gorm.io/gorm"
)

type WebhookLogRepository struct {
	db *gorm.DB
}

func NewWebhookLogRepository(db *gorm.DB) *WebhookLogRepository {
	return &WebhookLogRepository{db: db}
}

func (r *WebhookLogRepository) Create(ctx context.Context, log *model.WebhookLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *WebhookLogRepository) UpdateStatus(ctx context.Context, id interface{}, status string, errorMessage string) error {
	return r.db.WithContext(ctx).Model(&model.WebhookLog{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        status,
			"error_message": errorMessage,
		}).Error
}
