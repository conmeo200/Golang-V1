package web

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type CheckoutHandler struct{}

func NewCheckoutHandler() *CheckoutHandler {
	return &CheckoutHandler{}
}

func (h *CheckoutHandler) Index(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Thanh toán",
		"CartItems": []map[string]interface{}{
			{"ID": 1, "Name": "Áo Thun Nam Cổ Tròn Premium", "Price": "150.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Quantity": 2},
			{"ID": 2, "Name": "Giày Thể Thao Nam Nữ", "Price": "250.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Quantity": 1},
			{"ID": 4, "Name": "Tai Nghe Bluetooth", "Price": "350.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Quantity": 1},
		},
		"Total": "850.000đ",
	}

	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/pages/checkout/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *CheckoutHandler) Success(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Đặt hàng thành công",
	}

	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/pages/checkout/success.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}
