package main

import (
	"fmt"
	"log"
	"net/http"
)

// homePage handles requests to the root path "/"
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the API!")
	log.Println("Endpoint Hit: homePage")
}

// handleRequests sets up the routes and starts the server
func handleRequests() {
	// Register the homePage handler for the root path
	http.HandleFunc("/", homePage)

	// Start the HTTP server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}