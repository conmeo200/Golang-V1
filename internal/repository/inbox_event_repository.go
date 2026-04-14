package repository

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InboxEventRepository struct {
	db *gorm.DB
}

func NewInboxEventRepository(db *gorm.DB) *InboxEventRepository {
	return &InboxEventRepository{db: db}
}

func (r *InboxEventRepository) WithTx(tx *gorm.DB) *InboxEventRepository {
	return &InboxEventRepository{db: tx}
}

func (r *InboxEventRepository) Update(ctx context.Context, event *model.InboxEvent) error {
	return r.db.WithContext(ctx).Save(event).Error
}

func (r *InboxEventRepository) Create(ctx context.Context, event *model.InboxEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *InboxEventRepository) HasBeenProcessed(ctx context.Context, eventID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.InboxEvent{}).
		Where("event_id = ?", eventID).
		Count(&count).Error
	return count > 0, err
}
