package bootstrap

import (
	"github.com/conmeo200/Golang-V1/internal/domain/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.TokenBlacklist{},
		&model.Order{},
		&model.Role{},
		&model.Permission{},
		&model.TaxDeclaration{},
		&model.IncomeRecord{},
		&model.Dependent{},
		&model.Transaction{},
		&model.Payment{},
		&model.PaymentEvent{},
		&model.OutboxEvents{},
		&model.InboxEvent{},
		&model.WebhookLog{},
		&model.DeadLetterEvent{},
		// News Models
		&model.NewsUser{},
		&model.Category{},
		&model.Article{},
		&model.ArticleTrans{},
		&model.Tag{},
		&model.NewsComment{},
		&model.ArticleStats{},
		&model.ArticleViewLog{},
		&model.ArticleVersion{},
		// E-commerce Multi-vendor & Products
		&model.Shop{},
		&model.ShopAddress{},
		&model.ProductCategory{},
		&model.Product{},
		&model.ProductVariant{},
		&model.ProductImage{},
		// Inventory
		&model.Inventory{},
		&model.InventoryLog{},
		// Cart & Checkout
		&model.Cart{},
		&model.CartItem{},
		&model.ShippingAddress{},
		// Order Items (Order is already registered above)
		&model.OrderItem{},
		// Affiliate
		&model.AffiliateProfile{},
		&model.AffiliateLink{},
		&model.AffiliateCommission{},
	)
}
