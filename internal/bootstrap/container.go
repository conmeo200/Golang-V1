package bootstrap

import (
	"gorm.io/gorm"

	rabbitmqinfra "github.com/conmeo200/Golang-V1/internal/infrastructure/rabbitmq"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/stripe"

	authOut "github.com/conmeo200/Golang-V1/internal/module/auth/adapter/out"
	authSvc "github.com/conmeo200/Golang-V1/internal/module/auth"

	orderOut "github.com/conmeo200/Golang-V1/internal/module/order/adapter/out"
	orderSvc "github.com/conmeo200/Golang-V1/internal/module/order"

	paymentOut "github.com/conmeo200/Golang-V1/internal/module/payment/adapter/out"
	paymentSvc "github.com/conmeo200/Golang-V1/internal/module/payment"

	txOut "github.com/conmeo200/Golang-V1/internal/module/transaction/adapter/out"
	txSvc "github.com/conmeo200/Golang-V1/internal/module/transaction"

	userOut "github.com/conmeo200/Golang-V1/internal/module/user/adapter/out"
	userSvc "github.com/conmeo200/Golang-V1/internal/module/user"

	"github.com/conmeo200/Golang-V1/internal/module/auth/port"
	orderPort "github.com/conmeo200/Golang-V1/internal/module/order/port"
	paymentPort "github.com/conmeo200/Golang-V1/internal/module/payment/port"
	txPort "github.com/conmeo200/Golang-V1/internal/module/transaction/port"
	userPort "github.com/conmeo200/Golang-V1/internal/module/user/port"
)

type Container struct {
	Config   *Config
	DB       *gorm.DB
	RabbitMQ *rabbitmqinfra.RabbitMQ

	// repositories
	UserRepo        userPort.UserRepository
	TokenRepo       port.TokenRepository
	TransactionRepo txPort.TransactionRepository
	AuthRepo        port.AuthRepository
	OrderRepo       orderPort.OrderRepository
	PaymentRepo     paymentPort.PaymentRepository

	// services
	AuthService        port.AuthService
	OrderService       orderPort.OrderService
	PaymentService     paymentPort.PaymentService
	UserService        userPort.UserService
	TransactionService txPort.TransactionService

	// Infrastructure
	StripeService *stripe.StripeService

}

func BuildContainer(
	cfg *Config,
	db *gorm.DB,
	rabbit *rabbitmqinfra.RabbitMQ,
) *Container {

	c := &Container{
		Config:   cfg,
		DB:       db,
		RabbitMQ: rabbit,
	}

	// repositories
	c.UserRepo = userOut.NewUserRepository(db)
	c.TokenRepo = authOut.NewTokenRepository(db)
	c.TransactionRepo = txOut.NewTransactionRepository(db)
	c.AuthRepo = authOut.NewAuthRepository(db)
	c.OrderRepo = orderOut.NewOrderRepository(db)
	c.PaymentRepo = paymentOut.NewPaymentRepository(db)

	// outbox/inbox repos for payment
	outboxRepo := paymentOut.NewOutboxEventRepository(db)
	inboxRepo := paymentOut.NewInboxEventRepository(db)

	// services
	c.AuthService = authSvc.NewAuthService(
		c.AuthRepo,
		c.UserRepo,
		c.TokenRepo,
	)

	c.OrderService = orderSvc.NewOrderService(
		c.OrderRepo,
		nil,
	)

	c.PaymentService = paymentSvc.NewPaymentService(
		c.PaymentRepo,
		outboxRepo,
		inboxRepo,
	)

	c.UserService = userSvc.NewUserService(
		c.UserRepo,
	)

	c.TransactionService = txSvc.NewTransactionService(
		c.TransactionRepo,
	)

	c.StripeService = stripe.NewStripeService(
		cfg.StripeSecretKey,
		cfg.StripeWebhookSecret,
	)

	return c
}
