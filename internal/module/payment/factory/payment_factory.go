package factory


import (
	"fmt"

	"github.com/conmeo200/Golang-V1/internal/module/payment/ports"
)

type PaymentFactory struct {
	stripeProvider ports.PaymentProvider
}

func NewPaymentFactory(
	stripeProvider ports.PaymentProvider,
) *PaymentFactory {

	return &PaymentFactory{
		stripeProvider: stripeProvider,
	}
}

func (f *PaymentFactory) GetProvider(method string) (ports.PaymentProvider, error) {
	switch method {
	case ports.PaymentMethodStripe:
		return f.stripeProvider, nil
	default:
		return nil, fmt.Errorf("unsupported payment method: %s", method)
	}
}