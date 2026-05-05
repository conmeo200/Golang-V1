package api

import (
	//"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/conmeo200/Golang-V1/internal/infrastructure/logger"
	"github.com/conmeo200/Golang-V1/internal/module/payment"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/stripe"
	"github.com/google/uuid"
	stripego "github.com/stripe/stripe-go/v78"
)

type WebhookHandler struct {
	paymentService payment.PaymentServiceInterface
	stripeService  *stripe.StripeService
}

func NewWebhookHandler(paymentService payment.PaymentServiceInterface, stripeService *stripe.StripeService) *WebhookHandler {
	return &WebhookHandler{
		paymentService: paymentService,
		stripeService:  stripeService,
	}
}


func (h *WebhookHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)

	fmt.Printf("Stripe Webhook Payload:")
	fmt.Println("Server time:", time.Now())
	

	if err != nil {
		logger.ErrorLogger.Printf("Error reading stripe webhook body: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	sigHeader := r.Header.Get("Stripe-Signature")
	event, err := h.stripeService.ConstructEvent(payload, sigHeader)

	if err != nil {
		logger.ErrorLogger.Printf("Error verifying stripe webhook signature: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Extract ID and dispatch to generic handler
	fmt.Printf("Stripe Webhook Event: %s\n", event.Type)

	var intentID string
	var orderIDStr string
	var rawData map[string]interface{}

	err = json.Unmarshal(event.Data.Raw, &rawData)
	if err != nil {
		logger.ErrorLogger.Printf("Error parsing stripe raw data: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "payment_intent.succeeded", "payment_intent.payment_failed":
		var pi stripego.PaymentIntent
		json.Unmarshal(event.Data.Raw, &pi)
		intentID = pi.ID
		orderIDStr = pi.Metadata["order_id"]
		fmt.Printf("Stripe Webhook PaymentIntent: %v\n", pi)

	case "charge.refunded":
		var charge stripego.Charge
		json.Unmarshal(event.Data.Raw, &charge)
		// For charges, the PaymentIntent ID is often embedded
		if charge.PaymentIntent != nil {
			intentID = charge.PaymentIntent.ID
		}
		fmt.Printf("Stripe Webhook Charge: %v\n", charge)

	default:
		logger.StripeLogger.Printf("Unhandled stripe event type: %s", event.Type)
		w.WriteHeader(http.StatusOK)
		return
	}

	if intentID == "" {
		logger.StripeLogger.Printf("No intent ID found in event %s", event.Type)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Use intentID or event ID to ensure idempotency
	eventID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(event.ID))

	// Send payload to outbox if needed
	fmt.Printf("Stripe Webhook Event ID: %s\n", eventID)

	outboxPayload := map[string]interface{}{}
	if orderIDStr != "" {
		outboxPayload["order_id"] = orderIDStr
	}

	err = h.paymentService.HandleWebhookEvent(r.Context(), intentID, string(event.Type), eventID, outboxPayload)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to process payment from stripe webhook: %v", err)
	} else {
		logger.StripeLogger.Printf("🚀 Webhook processed | Event: %s | Intent: %s", event.Type, intentID)
	}

	fmt.Printf("Stripe Webhook Done")
	w.WriteHeader(http.StatusOK)
}
