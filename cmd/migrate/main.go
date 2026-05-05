package main

import (
	"log"

	"github.com/conmeo200/Golang-V1/internal/bootstrap"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/logger"
)

func main() {
	logger.Init()
	logger.AppLogger.Println("Initializing migration...")

	cfg := bootstrap.Load()
	db := bootstrap.InitDatabase(cfg)

	logger.AppLogger.Println("Running migrations...")
	if err := bootstrap.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	logger.AppLogger.Println("Migrations completed successfully.")
}
