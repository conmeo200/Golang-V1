package router

import (
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/handler"
	"github.com/conmeo200/Golang-V1/internal/help"
	"github.com/conmeo200/Golang-V1/internal/middleware"
)

func RegisterAuthRoutes(mux *http.ServeMux, h *handler.AuthHandler) {

	// Unprotected by JWT (but protected by API Key)
	mux.Handle("/api/v1/register", middleware.RequireAPIKey(help.Method(http.MethodPost, h.RegisterHandler)))
	mux.Handle("/api/v1/login", middleware.RequireAPIKey(help.Method(http.MethodPost, h.LoginHandler)))
	mux.Handle("/api/v1/forgot-password", middleware.RequireAPIKey(help.Method(http.MethodPost, h.ForgotPasswordHandler)))
	mux.Handle("/api/v1/refresh-token", middleware.RequireAPIKey(help.Method(http.MethodPost, h.RefreshTokenHandler)))

	// Protected by both API Key and JWT
	mux.Handle("/api/v1/logout", middleware.RequireAPIKey(middleware.JWTMiddleware(help.Method(http.MethodPost, h.LogoutHandler))))
	mux.Handle("/api/v1/change-password", middleware.RequireAPIKey(middleware.JWTMiddleware(help.Method(http.MethodPost, h.ChangePasswordHandler))))
	mux.Handle("/api/v1/revoke-token", middleware.RequireAPIKey(middleware.JWTMiddleware(help.Method(http.MethodPost, h.RevokeTokenHandler))))
}
