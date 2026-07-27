package bootstrap

import (
	"log"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/handler"
)

type APIApp struct {
	Container *Container
	Server    *http.Server
}

func NewAPIApp(container *Container) *APIApp {
	// Router (Gin)
	router := handler.InitRouter(container.AuthService)
	log.Println("✅ Web routes registered successfully")

	server := &http.Server{
		Addr:    ":" + container.Config.Port,
		Handler: router,
	}

	return &APIApp{
		Container: container,
		Server:    server,
	}
}

func (a *APIApp) Run() {
	log.Printf("🚀 Web Server starting on %s", a.Server.Addr)
	if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Server failed: %v", err)
	}
}

func (a *APIApp) Stop() {
	log.Println("🛑 Stopping Server...")
}
