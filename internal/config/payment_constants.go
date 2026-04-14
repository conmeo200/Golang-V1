package config

// Payment methods that the system supports
const (
	PaymentMethodPayPal = "paypal"
	PaymentMethodStripe = "stripe"
	PaymentMethodVNPay  = "vnpay"
	PaymentMethodCash   = "cash"
)

// Payment flow status constants
const (
	PaymentStatusPending    = "pending"
	PaymentStatusProcessing = "processing"
	PaymentStatusSuccess    = "success"
	PaymentStatusFailed     = "failed"
	PaymentStatusRefunded   = "refunded"
)

// Currency standard for the application
const DefaultCurrency = "USD"

// PaymentSettings structures hardcoded functional settings
type paymentSettings struct {
	SupportedMethods []string
	TaxRate          float64
	StripeFeeRate    float64
	PayPalFeeRate    float64
	GatewayTimeout   int // processing timeout in seconds
}

// PaymentConfig represents the global hardcoded configuration that does not change between environments
var PaymentConfig = paymentSettings{
	SupportedMethods: []string{PaymentMethodStripe, PaymentMethodPayPal, PaymentMethodVNPay},
	TaxRate:          0.00, // No tax default
	StripeFeeRate:    0.029, // Example: 2.9% + 30 cents (standard stripe fee)
	PayPalFeeRate:    0.034, // Example: 3.4% + 30 cents
	GatewayTimeout:   30,
}
