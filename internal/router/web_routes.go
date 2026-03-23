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

	// 2. Authentication Routes
	mux.HandleFunc("/login", help.Method(http.MethodGet, h.LoginPage))
	mux.HandleFunc("/register", help.Method(http.MethodGet, h.RegisterPage))

	// 3. Dashboard Routes
	mux.HandleFunc("/dashboard", help.Method(http.MethodGet, h.DashboardPage))
	mux.HandleFunc("/dashboard/users", help.Method(http.MethodGet, h.UserListPage))
	mux.HandleFunc("/dashboard/users/new", help.Method(http.MethodGet, h.UserFormPage))
	mux.HandleFunc("/dashboard/users/edit/1", help.Method(http.MethodGet, h.UserFormPage)) // Mocking specific edit route
	mux.HandleFunc("/dashboard/grades", help.Method(http.MethodGet, h.GradesPage))
	mux.HandleFunc("/dashboard/profile", help.Method(http.MethodGet, h.ProfilePage))
	mux.HandleFunc("/dashboard/analytics", help.Method(http.MethodGet, h.AnalyticsPage))
	mux.HandleFunc("/dashboard/sports", help.Method(http.MethodGet, h.SportsPage))
	mux.HandleFunc("/dashboard/sports/match/1", help.Method(http.MethodGet, h.MatchDetailPage))
	mux.HandleFunc("/dashboard/sports/standings", help.Method(http.MethodGet, h.StandingsPage))

	// 4. Presentation Routes (HTML templates)
	mux.HandleFunc("/news", help.Method(http.MethodGet, h.NewsPage))
    
    // Default route redirects to dashboard for now
    mux.HandleFunc("/", help.Method(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}))
}
