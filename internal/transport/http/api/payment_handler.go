package api

import (
	"encoding/json"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/bootstrap"
	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/conmeo200/Golang-V1/internal/module/payment"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/stripe"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService payment.PaymentServiceInterface
	stripeService  *stripe.StripeService
}

func NewPaymentHandler(paymentService payment.PaymentServiceInterface, stripeService *stripe.StripeService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		stripeService:  stripeService,
	}
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
		req.Currency = bootstrap.DefaultCurrency
	}

	factory := payment.NewPaymentFactory()
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

	intentID, _ := res["intent_id"].(string)

	// Convert orderID string to UUID
	orderUUID, err := uuid.Parse(req.OrderID)
	if err != nil {
		http.Error(w, `{"error": "Invalid order_id format. Must be a valid UUID."}`, http.StatusBadRequest)
		return
	}

	paymentUUid := uuid.New()
	// Insert PENDING payment record
	pendingPayment := &model.Payment{
		UUID:              paymentUUid,
		OrderID:           orderUUID,
		Amount:            req.Amount,
		Currency:          req.Currency,
		PaymentMethod:     req.PaymentMethod,
		Provider:          req.PaymentMethod,
		ProviderPaymentID: intentID,
		IdempotencyKey:    "payment_" + paymentUUid.String(),
	}

	if err := h.paymentService.CreatePendingPayment(r.Context(), pendingPayment); err != nil {
		http.Error(w, `{"error": "Failed to create pending payment"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *PaymentHandler) RefundPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	intentID, ok := req["intent_id"].(string)
	if !ok || intentID == "" {
		http.Error(w, `{"error": "intent_id is required"}`, http.StatusBadRequest)
		return
	}

	amountStr, _ := req["amount"].(float64) // Optional partial refund

	refundRes, err := h.stripeService.Refund(intentID, amountStr)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    refundRes.Status,
		"refund_id": refundRes.ID,
	})
}
