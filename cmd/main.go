package main

import (
	"log"
	"net/http"

	//"github.com/conmeo200/Golang-V1/internal/model"
	//"github.com/conmeo200/Golang-V1/database/seeder"
	"github.com/conmeo200/Golang-V1/internal/app"
	"github.com/conmeo200/Golang-V1/internal/auth"
	"github.com/conmeo200/Golang-V1/internal/database"

	//"github.com/conmeo200/Golang-V1/internal/handler"
	"github.com/conmeo200/Golang-V1/internal/router"
	//"github.com/conmeo200/Golang-V1/internal/service"
)

func main() {
	// 1. Load config

	// 2. Initialize dependencies

	// 3. Setup router

	// 4. Create server

	// Connect DB
	dbPostgres, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err)
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

	accessToken, refreshToken, err := auth.GenerateTokens("123")

	log.Println("GenerateTokens", accessToken, refreshToken, err)

	http.ListenAndServe(":8080", mux)

	log.Printf("Server starting on :%s\n", "8080")

	// Run server in goroutine

}
