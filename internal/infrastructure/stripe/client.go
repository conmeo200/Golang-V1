package stripe

import (
	"fmt"

	"github.com/conmeo200/Golang-V1/internal/bootstrap"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/logger"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"github.com/stripe/stripe-go/v78/refund"
	"github.com/stripe/stripe-go/v78/webhook"
)

type StripeService struct{
	webhookSecret string
}

func NewStripeService() *StripeService {
	// Initialize Stripe API credentials
	cfg := bootstrap.Load()
	if cfg.StripeSecretKey != "" {
		stripe.Key = cfg.StripeSecretKey
	} else {
		logger.ErrorLogger.Println("WARNING: Stripe Secret Key is missing in environment.")
	}
	
	return &StripeService{
		webhookSecret: cfg.StripeWebhookSecret,
	}
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

func (s *StripeService) ConstructEvent(payload []byte, sigHeader string) (stripe.Event, error) {
	return webhook.ConstructEventWithOptions(payload, sigHeader, s.webhookSecret, webhook.ConstructEventOptions{
		IgnoreTolerance:          true,
		IgnoreAPIVersionMismatch: true,
	})
}

// Refund thực hiện gọi API Stripe để hoàn tiền 
func (s *StripeService) Refund(intentID string, amount float64) (*stripe.Refund, error) {
	if stripe.Key == "" {
		return nil, fmt.Errorf("stripe api key is missing")
	}

	amountInCents := int64(amount * 100)
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(intentID),
	}
	
	// Nếu muốn refund một phần, thêm Amount. Nếu không truyền Amount thì Stripe sẽ refund toàn bộ.
	if amount > 0 {
		params.Amount = stripe.Int64(amountInCents)
	}

	r, err := refund.New(params)
	if err != nil {
		logger.StripeLogger.Printf("Stripe Refund Failed | Intent: %s | Error: %v", intentID, err)
		return nil, err
	}

	logger.StripeLogger.Printf("Stripe Refund Success | Intent: %s | RefundID: %s", intentID, r.ID)
	return r, nil
}
