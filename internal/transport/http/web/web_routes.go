package web

import (
	"fmt"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/transport/http/middleware"
	"github.com/conmeo200/Golang-V1/internal/transport/http/api"
	"github.com/conmeo200/Golang-V1/internal/help"
)

func RegisterWebRoutes(mux *http.ServeMux, clientHandler *ClientHandler, dashboardHandler *DashboardHandler, paymentHandler *api.PaymentHandler, authMiddleware *middleware.AuthMiddleware) {
	
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

	// 3. User Dashboard (Protected - Grouped with Sub-Mux)
	dashboardMux := http.NewServeMux()

	// Dashboard Home
	dashboardMux.Handle("/dashboard", help.Method(http.MethodGet, dashboardHandler.DashboardPage))
	
	// Users Management
	dashboardMux.Handle("/dashboard/users", authMiddleware.RequirePermission("users_read", help.Method(http.MethodGet, dashboardHandler.UserListPage)))
	dashboardMux.Handle("/dashboard/users/add", authMiddleware.RequirePermission("users_write", help.Method(http.MethodGet, dashboardHandler.UserFormPage)))
	dashboardMux.Handle("/dashboard/users/edit/1", authMiddleware.RequirePermission("users_write", help.Method(http.MethodGet, dashboardHandler.UserFormPage)))
	dashboardMux.Handle("/dashboard/users/detail/1", authMiddleware.RequirePermission("users_read", help.Method(http.MethodGet, dashboardHandler.UserDetailPage)))
	
	// Grades & Profile
	dashboardMux.Handle("/dashboard/grades", help.Method(http.MethodGet, dashboardHandler.GradesPage)) 
	dashboardMux.Handle("/dashboard/profile", help.Method(http.MethodGet, dashboardHandler.ProfilePage))

	// Roles Management
	dashboardMux.Handle("/dashboard/roles", authMiddleware.RequirePermission("roles_read", help.Method(http.MethodGet, dashboardHandler.RoleListPage)))
	dashboardMux.Handle("/dashboard/roles/new", authMiddleware.RequirePermission("roles_write", help.Method(http.MethodGet, dashboardHandler.RoleFormPage)))
	dashboardMux.Handle("/dashboard/roles/create", authMiddleware.RequirePermission("roles_write", help.Method(http.MethodPost, dashboardHandler.ProcessRoleCreate)))
	dashboardMux.Handle("/dashboard/roles/edit/", authMiddleware.RequirePermission("roles_write", help.Method(http.MethodGet, dashboardHandler.RoleFormPage)))
	dashboardMux.Handle("/dashboard/roles/update/", authMiddleware.RequirePermission("roles_write", help.Method(http.MethodPost, dashboardHandler.ProcessRoleUpdate)))
	dashboardMux.Handle("/dashboard/roles/delete/", authMiddleware.RequirePermission("roles_write", help.Method(http.MethodPost, dashboardHandler.ProcessRoleDelete)))
	dashboardMux.Handle("/dashboard/roles/detail/", authMiddleware.RequirePermission("roles_write", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodPost {
			dashboardHandler.RoleDetailPage(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Logs Management
	dashboardMux.Handle("/dashboard/logs", authMiddleware.RequirePermission("logs_read", help.Method(http.MethodGet, dashboardHandler.LogListPage)))
	dashboardMux.Handle("/dashboard/logs/detail/", authMiddleware.RequirePermission("logs_read", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			dashboardHandler.LogDetailPage(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Taxes Management
	dashboardMux.Handle("/dashboard/taxes", authMiddleware.RequirePermission("taxes_read", help.Method(http.MethodGet, dashboardHandler.TaxListPage)))
	dashboardMux.Handle("/dashboard/taxes/new", authMiddleware.RequirePermission("taxes_write", help.Method(http.MethodGet, dashboardHandler.TaxNewPage)))
	dashboardMux.Handle("/dashboard/taxes/create", authMiddleware.RequirePermission("taxes_write", help.Method(http.MethodPost, dashboardHandler.ProcessTaxDeclaration)))

	// Order Management
	dashboardMux.Handle("/dashboard/orders", authMiddleware.RequirePermission("orders_read", help.Method(http.MethodGet, dashboardHandler.OrderListPage)))
	dashboardMux.Handle("/dashboard/orders/new", authMiddleware.RequirePermission("orders_write", help.Method(http.MethodGet, dashboardHandler.OrderNewPage)))
	dashboardMux.Handle("/dashboard/orders/create", authMiddleware.RequirePermission("orders_write", help.Method(http.MethodPost, dashboardHandler.ProcessOrderCreate)))
	dashboardMux.Handle("/dashboard/orders/detail", authMiddleware.RequirePermission("orders_read", help.Method(http.MethodGet, dashboardHandler.OrderDetailPage)))
	dashboardMux.Handle("/dashboard/orders/update-status", authMiddleware.RequirePermission("orders_write", help.Method(http.MethodPost, dashboardHandler.UpdateOrderStatus)))

	// Payment Management
	dashboardMux.Handle("/dashboard/payments", authMiddleware.RequirePermission("payments_read", help.Method(http.MethodGet, dashboardHandler.PaymentListPage)))
	dashboardMux.Handle("/dashboard/payments/new", authMiddleware.RequirePermission("payments_write", help.Method(http.MethodGet, dashboardHandler.PaymentNewPage)))
	dashboardMux.Handle("/dashboard/payments/create", authMiddleware.RequirePermission("payments_write", help.Method(http.MethodPost, dashboardHandler.ProcessPayment)))
	dashboardMux.Handle("/dashboard/payments/detail", authMiddleware.RequirePermission("payments_read", help.Method(http.MethodGet, dashboardHandler.PaymentDetailPage)))
	
	dashboardMux.Handle("/dashboard/payment", authMiddleware.RequirePermission("payments_read", help.Method(http.MethodGet, dashboardHandler.PaymentListPage)))
	dashboardMux.Handle("/dashboard/payment/new", authMiddleware.RequirePermission("payments_write", help.Method(http.MethodGet, dashboardHandler.PaymentNewPage)))
	dashboardMux.Handle("/dashboard/payment/create", authMiddleware.RequirePermission("payments_write", help.Method(http.MethodPost, dashboardHandler.ProcessPayment)))
	dashboardMux.Handle("/dashboard/payment/detail", authMiddleware.RequirePermission("payments_read", help.Method(http.MethodGet, dashboardHandler.PaymentDetailPage)))
	dashboardMux.Handle("/dashboard/transactions/process-payment", authMiddleware.RequirePermission("payments_write", http.HandlerFunc(paymentHandler.ProcessPayment)))

	// News Management
	dashboardMux.Handle("/dashboard/news/articles", authMiddleware.RequirePermission("news_read", help.Method(http.MethodGet, dashboardHandler.ArticleListPage)))
	dashboardMux.Handle("/dashboard/news/articles/new", authMiddleware.RequirePermission("news_write", help.Method(http.MethodGet, dashboardHandler.ArticleNewPage)))
	dashboardMux.Handle("/dashboard/news/articles/create", authMiddleware.RequirePermission("news_write", help.Method(http.MethodPost, dashboardHandler.ProcessArticleCreate)))
	
	dashboardMux.Handle("/dashboard/news/categories", authMiddleware.RequirePermission("news_read", help.Method(http.MethodGet, dashboardHandler.CategoryListPage)))
	dashboardMux.Handle("/dashboard/news/categories/create", authMiddleware.RequirePermission("news_write", help.Method(http.MethodPost, dashboardHandler.ProcessCategoryCreate)))

	// Apply RequireAuth globally to the dashboard group
	mux.Handle("/dashboard/", authMiddleware.RequireAuth(dashboardMux))
	mux.Handle("/dashboard", authMiddleware.RequireAuth(dashboardMux))





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
