package seeder

import (
	"log"
	"time"

	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedOrdersAndTransactions(db *gorm.DB) error {
	var admin model.User
	if err := db.Where("email = ?", "admin@example.com").First(&admin).Error; err != nil {
		log.Printf("Admin user not found for seeding orders, skipping...")
		return nil
	}

	// 1. Seed Orders
	now := time.Now().Unix()
	orders := []model.Order{
		{
			UUID:           uuid.New(), // Fixed ID for seeder consistency if needed, but uuid.New() is fine
			UserID:         admin.ID,
			Amount:         250.50,
			Status:         "completed",
			PaymentStatus:  "paid",
			IdempotencyKey: "seed-order-1",
			CreatedAt:      now - 86400*2, // 2 days ago
		},
		{
			UUID:           uuid.New(),
			UserID:         admin.ID,
			Amount:         100.00,
			Status:         "pending",
			PaymentStatus:  "unpaid",
			IdempotencyKey: "seed-order-2",
			CreatedAt:      now - 86400, // 1 day ago
		},
	}

	for i := range orders {
		err := db.Where(model.Order{IdempotencyKey: orders[i].IdempotencyKey}).Attrs(orders[i]).FirstOrCreate(&orders[i]).Error
		if err != nil {
			return err
		}
	}

	return nil
}
