package handler

import (
	"github.com/gin-gonic/gin"
)

// InitRouter initializes the Gin router and registers all module routes
func InitRouter(webHandler *WebHandler) *gin.Engine {
	r := gin.Default()

	// Serve static files
	r.Static("/static", "web/static")

	// Web routes
	r.GET("/", webHandler.Home)

	// API routes group
	// api := r.Group("/api")
	// TODO: Inject module handlers here
	// e.g. orderHandler.RegisterRoutes(api.Group("/orders"))

	return r
}
