package persistence

import (
	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type deadLetterRepo struct {
	db *gorm.DB
}

func NewDeadLetterRepository(db *gorm.DB) model.DeadLetterRepo {
	return &deadLetterRepo{db: db}
}

func (r *deadLetterRepo) Create(event *model.DeadLetterEvent) error {
	return r.db.Create(event).Error
}

func (r *deadLetterRepo) GetByID(id uint) (*model.DeadLetterEvent, error) {
	var event model.DeadLetterEvent
	err := r.db.First(&event, id).Error
	return &event, err
}

func (r *deadLetterRepo) ListPending() ([]model.DeadLetterEvent, error) {
	var events []model.DeadLetterEvent
	err := r.db.Where("status = ?", "pending").Find(&events).Error
	return events, err
}

func (r *deadLetterRepo) MarkAsReplayed(id uint) error {
	return r.db.Model(&model.DeadLetterEvent{}).Where("id = ?", id).Update("status", "replayed").Error
}

func (r *deadLetterRepo) MarkAsResolved(id uint) error {
	return r.db.Model(&model.DeadLetterEvent{}).Where("id = ?", id).Update("status", "resolved").Error
}

func (r *deadLetterRepo) Exists(eventID uuid.UUID, queueName string) (bool, error) {
	var count int64
	err := r.db.Model(&model.DeadLetterEvent{}).
		Where("event_id = ? AND queue_name = ?", eventID, queueName).
		Count(&count).Error
	return count > 0, err
}
