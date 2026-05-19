package web

import (
	authmodule "github.com/conmeo200/Golang-V1/internal/module/auth"
	"github.com/conmeo200/Golang-V1/internal/module/news"
	"github.com/conmeo200/Golang-V1/internal/module/order"
	"github.com/conmeo200/Golang-V1/internal/module/payment"
	"github.com/conmeo200/Golang-V1/internal/module/role"
	"github.com/conmeo200/Golang-V1/internal/module/tax"
	"github.com/conmeo200/Golang-V1/internal/module/transaction"
	"github.com/conmeo200/Golang-V1/internal/module/user"
)

type DashboardHandler struct {
	authService  authmodule.AuthServiceInterface
	roleService  role.RoleServiceInterface
	taxService         tax.TaxServiceInterface
	orderService       order.OrderServiceInterface
	transactionService transaction.TransactionServiceInterface
	paymentService     payment.PaymentServiceInterface
	userService        user.UserServiceInterface
	newsService        *news.NewsService
}

func NewDashboardHandler(
	authService authmodule.AuthServiceInterface,
	roleService role.RoleServiceInterface,
	taxService tax.TaxServiceInterface,
	orderService order.OrderServiceInterface,
	transactionService transaction.TransactionServiceInterface,
	paymentService payment.PaymentServiceInterface,
	userService user.UserServiceInterface,
	newsService *news.NewsService,
) *DashboardHandler {
	return &DashboardHandler{
		authService:        authService,
		roleService:        roleService,
		taxService:         taxService,
		orderService:       orderService,
		transactionService: transactionService,
		paymentService:     paymentService,
		userService:        userService,
		newsService:        newsService,
	}
}
