package router

import (
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/app"
	api_router "github.com/conmeo200/Golang-V1/internal/router/api"
	web_router "github.com/conmeo200/Golang-V1/internal/router/web"
)

func RegisterRoutes(mux *http.ServeMux, app *app.App) {

	// Route Api
	api_router.RegisterAPIRoutes(mux, app.UserHandler, app.AuthHandler, app.OrderHandler)

	// Route Web
	web_router.RegisterWebRoutes(mux, app.ClientHandler, app.DashboardHandler)

}
