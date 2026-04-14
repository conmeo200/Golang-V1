package web

import (
	"fmt"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/handler/web/client"
	"github.com/conmeo200/Golang-V1/internal/handler/web/dashboard"
	"github.com/conmeo200/Golang-V1/internal/help"
)

func RegisterWebRoutes(mux *http.ServeMux, clientHandler *client.ClientHandler, dashboardHandler *dashboard.DashboardHandler) {
	
	// 1. Serve static files (CSS, JS, Images)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// 2. Authentication Pages
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			clientHandler.LoginPage(w, r)
		} else if r.Method == http.MethodPost {
			clientHandler.ProcessClientLogin(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/register", clientHandler.RegisterPage)
	mux.HandleFunc("/logout", clientHandler.Logout)

	mux.HandleFunc("/dashboard/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			dashboardHandler.DashboardLoginPage(w, r)
		} else if r.Method == http.MethodPost {
			dashboardHandler.ProcessDashboardLogin(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 3. User Dashboard (Protected)
	mux.Handle("/dashboard", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.DashboardPage)))
	mux.Handle("/dashboard/users", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.UserListPage)))
	mux.Handle("/dashboard/users/add", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.UserFormPage)))
	mux.Handle("/dashboard/users/edit/1", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.UserFormPage)))
	mux.Handle("/dashboard/grades", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.GradesPage)))
	mux.Handle("/dashboard/profile", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.ProfilePage)))
	mux.Handle("/dashboard/analytics", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.AnalyticsPage)))
	mux.Handle("/dashboard/sports", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, clientHandler.SportsPage)))
	mux.Handle("/dashboard/sports/match/1", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, clientHandler.MatchDetailPage)))
	mux.Handle("/dashboard/sports/standings", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, clientHandler.StandingsPage)))

	mux.Handle("/dashboard/users/detail/1", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.UserDetailPage)))

	// Roles Management
	mux.Handle("/dashboard/roles", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.RoleListPage)))
	mux.Handle("/dashboard/roles/detail/", dashboardHandler.DashboardMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodPost {
			dashboardHandler.RoleDetailPage(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Logs Management
	mux.Handle("/dashboard/logs", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.LogListPage)))
	mux.Handle("/dashboard/logs/detail/", dashboardHandler.DashboardMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			dashboardHandler.LogDetailPage(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Taxes Management
	mux.Handle("/dashboard/taxes", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.TaxListPage)))
	mux.Handle("/dashboard/taxes/new", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.TaxNewPage)))
	mux.Handle("/dashboard/taxes/create", dashboardHandler.DashboardMiddleware(help.Method(http.MethodPost, dashboardHandler.ProcessTaxDeclaration)))

	// Order Management
	mux.Handle("/dashboard/orders", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.OrderListPage)))
	mux.Handle("/dashboard/orders/detail", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.OrderDetailPage)))
	mux.Handle("/dashboard/orders/update-status", dashboardHandler.DashboardMiddleware(help.Method(http.MethodPost, dashboardHandler.UpdateOrderStatus)))

	// Payment Management (Support both plural and singular)
	mux.Handle("/dashboard/payments", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.PaymentListPage)))
	mux.Handle("/dashboard/payments/new", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.PaymentNewPage)))
	mux.Handle("/dashboard/payments/create", dashboardHandler.DashboardMiddleware(help.Method(http.MethodPost, dashboardHandler.ProcessPayment)))
	
	mux.Handle("/dashboard/payment", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.PaymentListPage)))
	mux.Handle("/dashboard/payment/new", dashboardHandler.DashboardMiddleware(help.Method(http.MethodGet, dashboardHandler.PaymentNewPage)))
	mux.Handle("/dashboard/payment/create", dashboardHandler.DashboardMiddleware(help.Method(http.MethodPost, dashboardHandler.ProcessPayment)))

	// 4. Presentation Routes (HTML templates)
	mux.HandleFunc("/", help.Method(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			fmt.Printf("🔍 [DEBUG] 404 Not Found at path: %s\n", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		clientHandler.HomePage(w, r)
	}))
}
