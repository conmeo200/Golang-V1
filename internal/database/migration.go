package database

import (
	"github.com/conmeo200/Golang-V1/internal/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.TokenBlacklist{},
		&model.Order{},
		&model.Role{},
		&model.Permission{},
		&model.TaxDeclaration{},
		&model.IncomeRecord{},
		&model.Dependent{},
		&model.Transaction{},
		&model.Payment{},
		&model.PaymentEvent{},
		&model.OutboxEvents{},
		&model.InboxEvent{},
		&model.WebhookLog{},
	)
}
