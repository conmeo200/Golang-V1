package router

import (
	"net/http"
	"github.com/conmeo200/Golang-V1/internal/handler"
)

func New(
	userHandler *handler.UserHandler,
	//orderHandler *handler.OrderHandler,
) *http.ServeMux {

	mux := http.NewServeMux()

	registerUserRoutes(mux, userHandler)
	//registerOrderRoutes(mux, orderHandler)

	return mux
}