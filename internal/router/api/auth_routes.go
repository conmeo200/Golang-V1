package api

import (
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/handler/api"
	"github.com/conmeo200/Golang-V1/internal/help"
	"github.com/conmeo200/Golang-V1/internal/middleware"
)

func RegisterAuthRoutes(mux *http.ServeMux, authHandler *api.AuthHandler) {

	// Unprotected by JWT (but protected by API Key)
	mux.Handle("/api/v1/register", middleware.RequireAPIKey(help.Method(http.MethodPost, authHandler.Register)))
	mux.Handle("/api/v1/login", middleware.RequireAPIKey(help.Method(http.MethodPost, authHandler.Login)))
	mux.Handle("/api/v1/forgot-password", middleware.RequireAPIKey(help.Method(http.MethodPost, authHandler.ForgotPassword)))
	mux.Handle("/api/v1/refresh-token", middleware.RequireAPIKey(help.Method(http.MethodPost, authHandler.RefreshToken)))

	// Protected by both API Key and JWT
	mux.Handle("/api/v1/logout", middleware.RequireAPIKey(middleware.JWTMiddleware(help.Method(http.MethodPost, authHandler.Logout))))
	mux.Handle("/api/v1/change-password", middleware.RequireAPIKey(middleware.JWTMiddleware(help.Method(http.MethodPost, authHandler.ChangePassword))))
	mux.Handle("/api/v1/revoke-token", middleware.RequireAPIKey(middleware.JWTMiddleware(help.Method(http.MethodPost, authHandler.RevokeToken))))
}
