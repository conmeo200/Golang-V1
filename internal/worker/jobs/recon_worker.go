package jobs

import (
	"context"
	"log"
	"time"

	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/persistence"
	"github.com/conmeo200/Golang-V1/internal/module/payment"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
)

type ReconciliationWorker struct {
	paymentRepo  persistence.PaymentRepo
	paymentSvc   payment.PaymentServiceInterface
}

func NewReconciliationWorker(paymentRepo persistence.PaymentRepo, paymentSvc payment.PaymentServiceInterface) *ReconciliationWorker {
	return &ReconciliationWorker{
		paymentRepo: paymentRepo,
		paymentSvc:  paymentSvc,
	}
}

func (w *ReconciliationWorker) Name() string {
	return "reconciliation_worker"
}

func (w *ReconciliationWorker) Start(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	log.Println("🚀 Reconciliation Worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Reconciliation Worker context cancelled...")
			return nil
		case <-ticker.C:
			w.runReconciliation(ctx)
		}
	}
}

func (w *ReconciliationWorker) Stop() error {
	return nil
}

func (w *ReconciliationWorker) runReconciliation(ctx context.Context) {
	log.Println("🔄 Running Stripe Reconciliation Job...")

	// Fetch payments that are PENDING and older than 15 minutes
	cutoff := time.Now().Add(-15 * time.Minute).Unix()

	// Need a custom query on PaymentRepo to fetch these
	var payments []model.Payment
	err := w.paymentRepo.DB().WithContext(ctx).
		Where("status = ? AND created_at <= ?", "PENDING", cutoff).
		Find(&payments).Error

	if err != nil {
		log.Printf("❌ Failed to fetch pending payments for reconciliation: %v", err)
		return
	}

	for _, p := range payments {
		if p.Provider == "stripe" && p.ProviderPaymentID != "" {
			w.reconcileStripePayment(ctx, &p)
		}
	}
}

func (w *ReconciliationWorker) reconcileStripePayment(ctx context.Context, p *model.Payment) {
	pi, err := paymentintent.Get(p.ProviderPaymentID, nil)
	if err != nil {
		log.Printf("⚠️ Failed to get payment intent from Stripe for payment %s: %v", p.UUID, err)
		return
	}

	// Fake an event ID
	eventID := p.UUID

	switch pi.Status {
	case stripe.PaymentIntentStatusSucceeded:
		log.Printf("✅ Reconciling payment %s to SUCCESS", p.UUID)
		// We call the webhook handler logic manually or abstract it
		err = w.paymentSvc.HandleWebhookEvent(ctx, p.ProviderPaymentID, "payment_intent.succeeded", eventID, nil)
	case stripe.PaymentIntentStatusCanceled:
		log.Printf("❌ Reconciling payment %s to FAILED", p.UUID)
		err = w.paymentSvc.HandleWebhookEvent(ctx, p.ProviderPaymentID, "payment_intent.payment_failed", eventID, nil)
	default:
		// Requires payment method, processing, etc -> still pending or failed eventually
		if time.Now().Unix()-p.CreatedAt > 24*3600 {
			// Mark as failed if pending for more than 24h
			log.Printf("❌ Payment %s stuck in %s for >24h, failing", p.UUID, pi.Status)
			err = w.paymentSvc.HandleWebhookEvent(ctx, p.ProviderPaymentID, "payment_intent.payment_failed", eventID, nil)
		}
	}

	if err != nil {
		log.Printf("❌ Error applying reconciliation for payment %s: %v", p.UUID, err)
	}
}
