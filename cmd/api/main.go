package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/conmeo200/Golang-V1/internal/app"
	"github.com/conmeo200/Golang-V1/internal/bootstrap"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/logger"
)

func main() {
	logger.Init()
	logger.AppLogger.Println("Starting API Server...")

	// 1. Centralized Bootstrap
	container, err := bootstrap.InitContainer()

	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	
	defer container.Close()

	// 2. Initialize API App
	apiApp := app.NewAPIApp(container)

	// 3. Run in Background
	go apiApp.Run()

	// 4. Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.AppLogger.Println("API Server shutting down...")
	apiApp.Stop()
}
