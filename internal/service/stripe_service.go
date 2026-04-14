package service

import (
	"fmt"

	"github.com/conmeo200/Golang-V1/internal/config"
	"github.com/conmeo200/Golang-V1/internal/logger"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

type StripeService struct{}

func NewStripeService() *StripeService {
	// Initialize Stripe API credentials
	cfg := config.Load()
	if cfg.StripeSecretKey != "" {
		stripe.Key = cfg.StripeSecretKey
	} else {
		logger.ErrorLogger.Println("WARNING: Stripe Secret Key is missing in environment.")
	}
	
	return &StripeService{}
}

func (s *StripeService) AuthorizePayment(amount float64, currency string, orderID string) (map[string]interface{}, error) {
	logger.StripeLogger.Printf("Initiating Stripe Authorization | Order: %s | Amount: %.2f %s", orderID, amount, currency)

	if stripe.Key == "" {
		logger.StripeLogger.Printf("ERROR: Stripe API Key is not set.")
		return nil, fmt.Errorf("internal payment gateway error") // Send generic error to client
	}

	// Convert abstract amount to integer cents format for Stripe API (e.g. $10.99 -> 1099)
	amountInCents := int64(amount * 100)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountInCents),
		Currency: stripe.String(currency),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Metadata: map[string]string{
			"order_id": orderID,
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		logger.StripeLogger.Printf("Stripe API Call Failed | Order: %s | Error: %v", orderID, err)
		return nil, fmt.Errorf("payment processor declined the request")
	}

	logger.StripeLogger.Printf("Stripe Intent Created | Order: %s | IntentID: %s", orderID, pi.ID)

	response := map[string]interface{}{
		"client_secret": pi.ClientSecret,
		"intent_id":     pi.ID,
		"status":        string(pi.Status),
	}
	return response, nil
}
