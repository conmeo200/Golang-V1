package web

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Index(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Trang chủ",
		"Products": []map[string]interface{}{
			{"ID": 1, "Name": "Áo Thun Nam Cổ Tròn Premium Cotton", "Price": "150.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Sold": "1,2k đã bán"},
			{"ID": 2, "Name": "Giày Thể Thao Nam Nữ Siêu Nhẹ", "Price": "250.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Sold": "5k đã bán"},
			{"ID": 3, "Name": "Balo Thời Trang Đi Học Chống Nước", "Price": "199.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Sold": "3,4k đã bán"},
			{"ID": 4, "Name": "Tai Nghe Bluetooth Không Dây", "Price": "350.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Sold": "10k đã bán"},
			{"ID": 5, "Name": "Đồng Hồ Thông Minh Màn Hình Cảm Ứng", "Price": "499.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Sold": "2,1k đã bán"},
			{"ID": 6, "Name": "Ốp Lưng Điện Thoại Trong Suốt", "Price": "25.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Sold": "50k đã bán"},
			{"ID": 7, "Name": "Kem Chống Nắng SPF 50+", "Price": "120.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Sold": "15k đã bán"},
			{"ID": 8, "Name": "Mũ Lưỡi Trai Nam Nữ Basic", "Price": "45.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Sold": "8,2k đã bán"},
		},
	}

	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/pages/home/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}
