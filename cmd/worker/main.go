package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/conmeo200/Golang-V1/internal/app"
	"github.com/conmeo200/Golang-V1/internal/bootstrap"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/logger"
)

func main() {
	logger.Init()
	logger.AppLogger.Println("Starting Worker Node...")

	// 1. Centralized Bootstrap (DB, RMQ, Config)
	container, err := bootstrap.InitContainer()

	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to initialize container: %v", err)
	}
	
	defer container.Close()

	// 2. Initialize Worker App
	workerApp := app.NewWorkerApp(container)
	workerApp.Run()

	// 3. Graceful Shutdown
	logger.AppLogger.Println("Workers started. Press Ctrl+C to stop.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.AppLogger.Println("Workers shutting down...")
	workerApp.Stop()
}
