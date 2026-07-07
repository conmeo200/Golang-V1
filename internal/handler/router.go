package handler

import (
	"github.com/conmeo200/Golang-V1/internal/handler/web"
	"github.com/gin-gonic/gin"
)

// InitRouter initializes the Gin router and registers all module routes
func InitRouter() *gin.Engine {
	r := gin.Default()

	// Serve static files
	r.Static("/static", "web/static")

	// Init Web Handlers
	homeHandler     := web.NewHomeHandler()
	productHandler  := web.NewProductHandler()
	cartHandler     := web.NewCartHandler()
	checkoutHandler := web.NewCheckoutHandler()
	accountHandler  := web.NewAccountHandler()
	authHandler     := web.NewAuthHandler()
	pageHandler     := web.NewPageHandler()

	// Web routes
	
	r.GET("/login", authHandler.Login)
	r.GET("/register", authHandler.Register)
	
	// Info pages
	r.GET("/about", pageHandler.About)
	r.GET("/contact", pageHandler.Contact)
	r.GET("/blog", pageHandler.Blog)
	r.GET("/promotions", pageHandler.Promotions)

	//Home
	r.GET("/", homeHandler.Index)
	r.GET("/product/:id", productHandler.Detail)
	r.GET("/category", productHandler.Category)

	r.GET("/account/user", accountHandler.User)
	r.GET("/account/user/orders", accountHandler.UserOrders)
	r.GET("/account/user/notifications", accountHandler.UserNotifications)
	r.GET("/account/user/addresses", accountHandler.UserAddresses)
	r.GET("/account/user/password", accountHandler.UserPassword)

	r.GET("/seller", accountHandler.Seller)
	r.GET("/seller/products", accountHandler.SellerProducts)
	r.GET("/seller/products/add", accountHandler.SellerProductAdd)
	r.POST("/seller/products/add", accountHandler.SellerProductAddPost)
	r.GET("/seller/orders", accountHandler.SellerOrders)
	r.GET("/seller/revenue", accountHandler.SellerRevenue)
	r.GET("/seller/settings", accountHandler.SellerSettings)
	
	r.GET("/cart", cartHandler.Index)
	r.GET("/checkout", checkoutHandler.Index)
	r.GET("/checkout/success", checkoutHandler.Success)
	
	return r
}
