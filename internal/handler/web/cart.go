package web

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type CartHandler struct{}

func NewCartHandler() *CartHandler {
	return &CartHandler{}
}

func (h *CartHandler) Index(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Giỏ hàng",
		"CartItems": []map[string]interface{}{
			{"ID": 1, "Name": "Áo Thun Nam Cổ Tròn Premium Cotton", "Price": "150.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Quantity": 2},
			{"ID": 2, "Name": "Giày Thể Thao Nam Nữ Siêu Nhẹ", "Price": "250.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Quantity": 1},
			{"ID": 4, "Name": "Tai Nghe Bluetooth Không Dây", "Price": "350.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Quantity": 1},
		},
		"Subtotal": "900.000đ",
		"Total":    "850.000đ",
	}

	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/pages/cart/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}
