package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/google/uuid"
)

type PaymentListPageData struct {
	Title      string
	ActiveMenu string
	Payments   []model.Payment
}

type PaymentFormPageData struct {
	Title      string
	ActiveMenu string
	Orders     []model.Order
	Errors     map[string]string
}

type PaymentDeletePageData struct {
	Title      string
	ActiveMenu string
}

func (h *DashboardHandler) PaymentPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/payment.html")
	if err != nil {
		log.Printf("Error parsing payment template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := DashboardData{
		Title:      "Payment",
		ActiveMenu: "payment",
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *DashboardHandler) PaymentListPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/payment_list.html")
	if err != nil {
		log.Printf("Error parsing payment list template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	payments, err := h.paymentService.ListAllTransactions(r.Context())
	if err != nil {
		log.Printf("Error listing payments: %v", err)
	}

	data := PaymentListPageData{
		Title:      "Payment Management",
		ActiveMenu: "payments",
		Payments:   payments,
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *DashboardHandler) PaymentNewPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/payment_form.html")
	if err != nil {
		log.Printf("Error parsing payment form template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	orders, err := h.orderService.ListAllOrders(r.Context())
	if err != nil {
		log.Printf("Error fetching orders for payment: %v", err)
	}

	data := PaymentFormPageData{
		Title:      "New Payment",
		ActiveMenu: "payments",
		Orders:     orders,
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *DashboardHandler) PaymentEditPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/payment_form.html")
	if err != nil {
		log.Printf("Error parsing payment form template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := PaymentFormPageData{
		Title:      "Edit Payment",
		ActiveMenu: "payments",
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *DashboardHandler) PaymentDeletePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/payment_delete.html")
	if err != nil {
		log.Printf("Error parsing payment delete template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := PaymentDeletePageData{
		Title:      "Delete Payment",
		ActiveMenu: "payments",
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *DashboardHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/dashboard/payments/new", http.StatusSeeOther)
		return
	}

	amountStr := r.FormValue("amount")
	amount, _ := strconv.ParseFloat(amountStr, 64)
	currency := r.FormValue("currency")
	orderIDStr := r.FormValue("order_id")
	orderID, _ := uuid.Parse(orderIDStr)

	paymentUUid := uuid.New()
	payment := &model.Payment{
		UUID:           paymentUUid,
		OrderID:        orderID,
		Amount:         amount,
		Currency:       currency,
		PaymentMethod:  "manual", // Simulating a manual create for mechanism testing
		Status:         "PENDING",
		IdempotencyKey: "manual_" + paymentUUid.String(),
	}

	err := h.paymentService.CreatePayment(r.Context(), payment)
	if err != nil {
		log.Printf("Error creating payment: %v", err)
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard/payments", http.StatusSeeOther)
}

func (h *DashboardHandler) PaymentDetailPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/payment_detail.html")
	if err != nil {
		log.Printf("Error parsing payment detail template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	paymentUUIDStr := r.URL.Query().Get("uuid")
	paymentUUID, err := uuid.Parse(paymentUUIDStr)
	if err != nil {
		http.Redirect(w, r, "/dashboard/payments?error=Invalid+Payment+ID", http.StatusSeeOther)
		return
	}

	payment, err := h.paymentService.GetPaymentByUUID(r.Context(), paymentUUID)
	if err != nil || payment == nil {
		http.Redirect(w, r, "/dashboard/payments?error=Payment+not+found", http.StatusSeeOther)
		return
	}

	order, _ := h.orderService.GetOrder(r.Context(), payment.OrderID)

	data := struct {
		Title      string
		ActiveMenu string
		Payment    *model.Payment
		Order      *model.Order
	}{
		Title:      "Payment Detail #" + payment.UUID.String()[:8],
		ActiveMenu: "payments",
		Payment:    payment,
		Order:      order,
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}
