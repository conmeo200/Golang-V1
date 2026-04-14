package repository

import (
	"context"
	"time"

	"github.com/conmeo200/Golang-V1/internal/model"
	"gorm.io/gorm"
)

type OutboxEventRepository struct {
	db *gorm.DB
}

func NewOutboxEventRepository(db *gorm.DB) *OutboxEventRepository {
	return &OutboxEventRepository{db: db}
}

func (r *OutboxEventRepository) WithTx(tx *gorm.DB) *OutboxEventRepository {
	return &OutboxEventRepository{db: tx}
}

func (r *OutboxEventRepository) Create(ctx context.Context, event *model.OutboxEvents) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *OutboxEventRepository) Update(ctx context.Context, event *model.OutboxEvents) error {
	return r.db.WithContext(ctx).Save(event).Error
}

func (r *OutboxEventRepository) FetchPending(ctx context.Context, limit int) ([]model.OutboxEvents, error) {
	var events []model.OutboxEvents
	now := time.Now().Unix()
	err := r.db.WithContext(ctx).
		Where("status = ? AND next_retry_at <= ?", "PENDING", now).
		Order("created_at asc").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (r *OutboxEventRepository) MarkAsPublished(ctx context.Context, id interface{}, sentAt int64) error {
	return r.db.WithContext(ctx).Model(&model.OutboxEvents{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":  "SENT",
			"sent_at": sentAt,
		}).Error
}
