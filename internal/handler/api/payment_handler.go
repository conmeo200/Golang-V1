package api

import (
	"encoding/json"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/config"
	"github.com/conmeo200/Golang-V1/internal/service"
)

type PaymentHandler struct{}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

// PaymentRequest payload from mobile client
type PaymentRequest struct {
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	OrderID       string  `json:"order_id"`
}

func (h *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validation
	if req.PaymentMethod == "" {
		http.Error(w, `{"error": "payment_method is required"}`, http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		http.Error(w, `{"error": "amount must be greater than 0"}`, http.StatusBadRequest)
		return
	}
	if req.OrderID == "" {
		http.Error(w, `{"error": "order_id is required"}`, http.StatusBadRequest)
		return
	}

	if req.Currency == "" {
		req.Currency = config.DefaultCurrency
	}

	factory := service.NewPaymentFactory()
	provider, err := factory.GetProvider(req.PaymentMethod)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	res, err := provider.AuthorizePayment(req.Amount, req.Currency, req.OrderID)
	if err != nil {
		http.Error(w, `{"error": "Failed to process payment, please try again."}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
