package dashboard

import (
	"html/template"
	"log"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/google/uuid"
)

// OrderListPageData holds data for the order list page
type OrderListPageData struct {
	Title      string
	ActiveMenu string
	Orders     []model.Order
	Stats      OrderStats
	Error      string
	Success    string
}

type OrderStats struct {
	TotalOrders     int
	PendingOrders   int
	CompletedOrders int
	TotalRevenue    float64
}

// OrderListPage renders the list of all orders
func (h *DashboardHandler) OrderListPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/order_list.html")
	if err != nil {
		log.Printf("Error parsing order list template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	orders, err := h.orderService.ListAllOrders(r.Context())
	if err != nil {
		log.Printf("Error listing orders: %v", err)
	}

	stats := OrderStats{}
	stats.TotalOrders = len(orders)
	for _, o := range orders {
		if o.Status == "pending" {
			stats.PendingOrders++
		} else if o.Status == "completed" {
			stats.CompletedOrders++
			stats.TotalRevenue += o.Amount
		}
	}

	data := OrderListPageData{
		Title:      "Order Management",
		ActiveMenu: "orders",
		Orders:     orders,
		Stats:      stats,
		Success:    r.URL.Query().Get("success"),
		Error:      r.URL.Query().Get("error"),
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

// UpdateOrderStatus handles quick status updates from the dashboard
func (h *DashboardHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	// ... (content omitted for brevity but keeping the function)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderUUIDStr := r.FormValue("uuid")
	status := r.FormValue("status")
	paymentStatus := r.FormValue("payment_status")

	orderUUID, err := uuid.Parse(orderUUIDStr)
	if err != nil {
		http.Redirect(w, r, "/dashboard/orders?error=Invalid+Order+ID", http.StatusSeeOther)
		return
	}

	err = h.orderService.UpdateOrderStatus(r.Context(), orderUUID, status, paymentStatus)
	if err != nil {
		log.Printf("Error updating order status: %v", err)
		http.Redirect(w, r, "/dashboard/orders?error=Update+failed", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/dashboard/orders?success=Order+updated+successfully", http.StatusSeeOther)
}

// OrderDetailPage renders the detailed view of an order
func (h *DashboardHandler) OrderDetailPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/order_detail.html")
	if err != nil {
		log.Printf("Error parsing order detail template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	orderUUIDStr := r.URL.Query().Get("uuid")
	if orderUUIDStr == "" {
		// Try from path if using a custom router, but here using URL query for simplicity
		// consistent with how RoleDetailPage/LogDetailPage might work
	}

	orderUUID, err := uuid.Parse(orderUUIDStr)
	if err != nil {
		http.Redirect(w, r, "/dashboard/orders?error=Invalid+Order+ID", http.StatusSeeOther)
		return
	}

	order, err := h.orderService.GetOrder(r.Context(), orderUUID)
	if err != nil || order == nil {
		http.Redirect(w, r, "/dashboard/orders?error=Order+not+found", http.StatusSeeOther)
		return
	}

	transactions, _ := h.transactionService.GetTransactionsByOrderID(r.Context(), order.UUID)

	data := struct {
		Title        string
		ActiveMenu   string
		Order        *model.Order
		Transactions []model.Transaction
	}{
		Title:        "Order Detail #" + order.UUID.String()[:8],
		ActiveMenu:   "orders",
		Order:        order,
		Transactions: transactions,
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

// Helper to check if a render method exists, if not I'll implement it or use standard template execution
// Looking at tax.go, it uses tmpl.ExecuteTemplate. I'll stick to that or check if I should add renderDashboard.
