package main

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/conmeo200/Golang-V1/internal/app"
	"github.com/conmeo200/Golang-V1/internal/config"
	"github.com/conmeo200/Golang-V1/internal/database"
	"github.com/conmeo200/Golang-V1/internal/logger"
	"github.com/conmeo200/Golang-V1/internal/queue/rabbitmq"
	"github.com/conmeo200/Golang-V1/internal/worker"
)

func main() {
	cfg := config.Load()
	logger.Init()

	dbPostgres, err := database.NewPostgres(cfg)
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to connect to database: %v", err)
	}
	
	rabbitMQ, err := rabbitmq.NewRabbitMQ(cfg)
	if err != nil {
		logger.ErrorLogger.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	myApp := app.NewApp(dbPostgres, rabbitMQ)
    
    // Start Workers via Manager
    workerMgr := worker.NewManager(rabbitMQ, myApp.OrderService)
    workerMgr.Start()

	logger.AppLogger.Println("Workers started. Press Ctrl+C to stop.")
    
    // Graceful Shutdown signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    
    logger.AppLogger.Println("Workers shutting down...")
}
