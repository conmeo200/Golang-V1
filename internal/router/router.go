package router

import (
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/app"
)

func RegisterRoutes(mux *http.ServeMux, app *app.App) {

	RegisterUserRoutes(mux, app.UserHandler)
	RegisterAuthRoutes(mux, app.AuthHandler)

}