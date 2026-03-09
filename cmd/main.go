package main

import (
	// "context"
	// "net/http"
	// "os"
	// "os/signal"
	// "syscall"
	// "time"
	"log"
	"net/http"

	//"github.com/conmeo200/Golang-V1/internal/config"
	//"github.com/conmeo200/Golang-V1/internal/handler"
	//"github.com/conmeo200/Golang-V1/internal/repository"
	//"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/database"
	"github.com/conmeo200/Golang-V1/internal/handler"
	"github.com/conmeo200/Golang-V1/internal/router"
	"github.com/conmeo200/Golang-V1/internal/service"
)

func main() {
	// 1. Load config

	// 2. Initialize dependencies

	// 3. Setup router


	// 4. Create server

	dbPostgres, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	//Service, Handler and Route User
	userService := service.NewUserService(dbPostgres)
	userHandler := handler.NewUserHandler(userService)
	router := router.New(userHandler)

	log.Printf("Server starting on :%s\n", "8080")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}

	// Run server in goroutine
	
}