package main

import (
	"net/http"

	//"github.com/conmeo200/Golang-V1/internal/model"
	//"github.com/conmeo200/Golang-V1/database/seeder"
	"github.com/conmeo200/Golang-V1/internal/app"
	"github.com/conmeo200/Golang-V1/internal/config"
	"github.com/conmeo200/Golang-V1/internal/database"
	"github.com/conmeo200/Golang-V1/internal/logger"
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
	// dbPostgres.Migrator().DropTable(&model.User{})
	// err = dbPostgres.AutoMigrate(&model.User{})
	// err = dbPostgres.AutoMigrate(&model.TokenBlacklist{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Migration successfuly!")

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
