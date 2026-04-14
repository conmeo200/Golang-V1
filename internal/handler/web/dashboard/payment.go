package dashboard

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/conmeo200/Golang-V1/internal/model"
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

	data := PaymentFormPageData{
		Title:      "New Payment",
		ActiveMenu: "payments",
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

	payment := &model.Payment{
		UUID:          uuid.New(),
		OrderID:       orderID,
		Amount:        amount,
		Currency:      currency,
		PaymentMethod: "manual", // Simulating a manual create for mechanism testing
		Status:        "PENDING",
	}

	err := h.paymentService.CreatePayment(r.Context(), payment)
	if err != nil {
		log.Printf("Error creating payment: %v", err)
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard/payments", http.StatusSeeOther)
}