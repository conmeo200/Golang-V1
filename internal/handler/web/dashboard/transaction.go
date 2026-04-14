package dashboard

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/config"
	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/service"
)

type TransactionListPageData struct {
	Title        string
	ActiveMenu   string
	Transactions []model.Transaction
	Error        string
	Success      string
}

func (h *DashboardHandler) TransactionListPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/transaction_list.html")
	if err != nil {
		log.Printf("Error parsing transaction list template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	transactions, err := h.transactionService.ListAllTransactions(r.Context())
	if err != nil {
		log.Printf("Error listing transactions: %v", err)
	}

	data := TransactionListPageData{
		Title:        "Transaction Management",
		ActiveMenu:   "transactions",
		Transactions: transactions,
		Success:      r.URL.Query().Get("success"),
		Error:        r.URL.Query().Get("error"),
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

// PaymentRequest structure passed from the checkout UI layer via AJAX
type PaymentRequest struct {
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	OrderID       string  `json:"order_id"`
}

// ProcessPayment is the centralized REST endpoint that dispatches to the correct payment provider
func (h *DashboardHandler) EndPointProcessPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Currency == "" {
		req.Currency = config.DefaultCurrency
	}

	// 1. Utilize factory to dynamically choose provider (Stripe, Paypal, VNPay...)
	factory := service.NewPaymentFactory()
	provider, err := factory.GetProvider(req.PaymentMethod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 2. Delegate the actual processing logic
	res, err := provider.AuthorizePayment(req.Amount, req.Currency, req.OrderID)
	if err != nil {
		// Securely return generic error, since actual specific SDK logic is logged internally
		http.Error(w, `{"error": "Failed to authorize payment"}`, http.StatusInternalServerError)
		return
	}

	// 3. Return provider-specific tokens (like Stripe Client Secret) to UI
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
