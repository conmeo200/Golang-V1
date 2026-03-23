package main

import (
	//"fmt"
	"net/http"

	//"github.com/conmeo200/Golang-V1/internal/model"
	//"github.com/conmeo200/Golang-V1/database/seeder"
	"github.com/conmeo200/Golang-V1/internal/app"
	"github.com/conmeo200/Golang-V1/internal/config"
	"github.com/conmeo200/Golang-V1/internal/database"
	"github.com/conmeo200/Golang-V1/internal/logger"
	"github.com/conmeo200/Golang-V1/internal/queue/rabbitmq"
	consumer "github.com/conmeo200/Golang-V1/internal/queue/rabbitmq/consumers"

	//"github.com/conmeo200/Golang-V1/internal/queue/consumer"

	//"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/router"
)

func main() {
	// 1. Load config
	cfg := config.Load()
	// 2. Initialize custom loggers
	logger.Init()

	dbPostgres, err := database.NewPostgres(cfg)
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to connect to database: %v", err)
	}
	
	// 3. Init RabbitMQ
	rabbitMQ, err := rabbitmq.NewRabbitMQ(cfg)
	if err != nil {
		logger.ErrorLogger.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	//defer rabbitMQ.Close()

	// 5. Initialize App (Service, Handler and Route User)
	myApp := app.NewApp(dbPostgres, rabbitMQ)

	// 6. Start Consumer
	orderConsumer := consumer.NewOrderConsumer(rabbitMQ, myApp.OrderService)
	go orderConsumer.StartOrder()

	logger.AppLogger.Println("Starting application...")

	// Run Migration...
	// ... (Migration commented out)

	mux := http.NewServeMux()
	router.RegisterRoutes(mux, myApp)

	logger.AppLogger.Printf("Server starting on :%s\n", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, mux)
}
