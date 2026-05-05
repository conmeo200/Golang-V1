package main

import (
	"log"

	"github.com/conmeo200/Golang-V1/internal/bootstrap"
	"github.com/conmeo200/Golang-V1/internal/database/seeder"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/logger"
)

func main() {
	logger.Init()
	logger.AppLogger.Println("Initializing seeder...")

	cfg := bootstrap.Load()
	db := bootstrap.InitDatabase(cfg)

	logger.AppLogger.Println("Running seeders...")
	if err := seeder.Seed(db); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}
	logger.AppLogger.Println("Seeding completed successfully.")
}
