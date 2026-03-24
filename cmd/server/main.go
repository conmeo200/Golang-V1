package main

import (
	"net/http"
	"github.com/conmeo200/Golang-V1/internal/app"
	"github.com/conmeo200/Golang-V1/internal/config"
	"github.com/conmeo200/Golang-V1/internal/database"
	"github.com/conmeo200/Golang-V1/internal/logger"
	"github.com/conmeo200/Golang-V1/internal/queue/rabbitmq"
	"github.com/conmeo200/Golang-V1/internal/router"
)

func main() {
	// 1. Load config
	cfg := config.Load()
	// 2. Initialize custom loggers
	logger.Init()

	// 3. Init Database
	dbPostgres, err := database.NewPostgres(cfg)
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to connect to database: %v", err)
	}
	
	// 4. Init RabbitMQ
	rabbitMQ, err := rabbitmq.NewRabbitMQ(cfg)
	if err != nil {
		logger.ErrorLogger.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	// 5. Initialize App
	myApp := app.NewApp(dbPostgres, rabbitMQ)

	// 6. Register Routes
	mux := http.NewServeMux()
	router.RegisterRoutes(mux, myApp)

	logger.AppLogger.Printf("Server starting on :%s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		logger.ErrorLogger.Fatalf("server failed: %v", err)
	}
}
