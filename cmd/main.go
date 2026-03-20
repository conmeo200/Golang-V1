package main

import (
	//"log"
	//"fmt"
	"net/http"

	//"github.com/conmeo200/Golang-V1/internal/model"
	//"github.com/conmeo200/Golang-V1/database/seeder"
	"github.com/conmeo200/Golang-V1/internal/app"
	"github.com/conmeo200/Golang-V1/internal/config"
	"github.com/conmeo200/Golang-V1/internal/database"
	"github.com/conmeo200/Golang-V1/internal/logger"

	//"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/router"
)

func main() {
	// 1. Load config
	cfg := config.Load()

	// 2. Initialize dependencies
	// Initialize custom loggers
	logger.Init()

	logger.AppLogger.Println("Starting application...")

	dbPostgres, err := database.NewPostgres(cfg)
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to connect to database: %v", err)
	}

	// Run Migration
	// err = dbPostgres.AutoMigrate(
	// 	// &model.User{},
	// 	// &model.TokenBlacklist{},
	// 	&model.Order{},
	// )
	// if err != nil {
	// 	logger.ErrorLogger.Fatalf("Migration failed: %v", err)
	// }

	// logger.AppLogger.Println("Migration successful!")

	//Run Seeder
	// err = seeder.SeedUsers(dbPostgres)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Seeder successfuly!")

	//Service, Handler and Route User
	app := app.NewApp(dbPostgres)

	mux := http.NewServeMux()

	router.RegisterRoutes(mux, app)

	//accessToken, refreshToken, err := auth.GenerateTokens("123")

	//log.Println("GenerateTokens", accessToken, refreshToken, err)

	logger.AppLogger.Printf("Server starting on :%s\n", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, mux)
}
