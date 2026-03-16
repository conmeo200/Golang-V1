package router

import (
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/handler"
	"github.com/conmeo200/Golang-V1/internal/help"
)

func RegisterWebRoutes(mux *http.ServeMux, h *handler.WebHandler) {
	
	// 1. Serve static files (CSS, JS, Images)
	// StripPrefix removes the "/static/" prefix before searching the file in "web/static" directory
	staticFileServer := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

	// 2. Presentation Routes (HTML templates)
	// You can add middleware.WebAuthMiddleware(...) to protect this route if you want
	// Right now News page is configured to be public. Uncomment below if you want it protected.
	// mux.Handle("/news", middleware.WebAuthMiddleware(help.Method(http.MethodGet, h.NewsPage)))
	mux.HandleFunc("/news", help.Method(http.MethodGet, h.NewsPage))
    
    // Default route redirects to news
    mux.HandleFunc("/", help.Method(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/news", http.StatusTemporaryRedirect)
	}))
}
