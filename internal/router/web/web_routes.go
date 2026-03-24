package web

import (
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/handler/web"
	"github.com/conmeo200/Golang-V1/internal/help"
)

func RegisterWebRoutes(mux *http.ServeMux, webHandler *web.WebHandler) {
	
	// 1. Serve static files (CSS, JS, Images)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// 2. Authentication Pages
	mux.HandleFunc("/login", webHandler.LoginPage)
	mux.HandleFunc("/register", webHandler.RegisterPage)

	// 3. User Dashboard
	mux.HandleFunc("/dashboard", help.Method(http.MethodGet, webHandler.DashboardPage))
	mux.HandleFunc("/dashboard/users", help.Method(http.MethodGet, webHandler.UserListPage))
	mux.HandleFunc("/dashboard/users/add", help.Method(http.MethodGet, webHandler.UserFormPage))
	mux.HandleFunc("/dashboard/users/edit/1", help.Method(http.MethodGet, webHandler.UserFormPage))
	mux.HandleFunc("/dashboard/grades", help.Method(http.MethodGet, webHandler.GradesPage))
	mux.HandleFunc("/dashboard/profile", help.Method(http.MethodGet, webHandler.ProfilePage))
	mux.HandleFunc("/dashboard/analytics", help.Method(http.MethodGet, webHandler.AnalyticsPage))
	mux.HandleFunc("/dashboard/sports", help.Method(http.MethodGet, webHandler.SportsPage))
	mux.HandleFunc("/dashboard/sports/match/1", help.Method(http.MethodGet, webHandler.MatchDetailPage))
	mux.HandleFunc("/dashboard/sports/standings", help.Method(http.MethodGet, webHandler.StandingsPage))

	// 4. Presentation Routes (HTML templates)
	mux.HandleFunc("/news", help.Method(http.MethodGet, webHandler.NewsPage))
    
    // Default route redirects to dashboard for now
    mux.HandleFunc("/", help.Method(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}))
}
