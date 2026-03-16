package router

import (
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/handler"
	"github.com/conmeo200/Golang-V1/internal/help"
	"github.com/conmeo200/Golang-V1/internal/middleware"
)

func RegisterAuthRoutes(mux *http.ServeMux, h *handler.AuthHandler) {

	mux.HandleFunc("/register", help.Method(http.MethodPost, h.RegisterHandler))
	mux.HandleFunc("/login", help.Method(http.MethodPost, h.LoginHandler))
	mux.HandleFunc("/forgot-password", help.Method(http.MethodPost, h.ForgotPasswordHandler))
	mux.HandleFunc("/refresh-token", help.Method(http.MethodPost, h.RefreshTokenHandler))

	mux.Handle("/logout", middleware.JWTMiddleware(help.Method(http.MethodPost, h.LogoutHandler)))
	mux.Handle("/change-password", middleware.JWTMiddleware(help.Method(http.MethodPost, h.ChangePasswordHandler)))
	mux.Handle("/revoke-token", middleware.JWTMiddleware(help.Method(http.MethodPost, h.RevokeTokenHandler)))
}
