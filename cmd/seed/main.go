package main

import (
	"github.com/conmeo200/Golang-V1/internal/config"
	"github.com/conmeo200/Golang-V1/internal/database"
	"github.com/conmeo200/Golang-V1/internal/database/seeder"
	"github.com/conmeo200/Golang-V1/internal/logger"
)

func main() {
	cfg := config.Load()
	logger.Init()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to connect to database: %v", err)
	}

	logger.AppLogger.Println("Running seeders...")
	if err := seeder.Seed(db); err != nil {
		logger.ErrorLogger.Fatalf("Seeding failed: %v", err)
	}
	logger.AppLogger.Println("Seeding completed successfully.")
}
