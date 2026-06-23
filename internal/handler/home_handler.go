package handler

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type WebHandler struct {
	tmpl *template.Template
}

func NewWebHandler() *WebHandler {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	return &WebHandler{tmpl: tmpl}
}

func (h *WebHandler) Home(c *gin.Context) {
	// Example product data for the view
	data := map[string]interface{}{
		"Title": "Shopee Clone - Mua sắm trực tuyến",
		"Products": []map[string]interface{}{
			{"Name": "Áo Thun Nam Cổ Tròn", "Price": "150.000đ", "Image": "https://down-vn.img.susercontent.com/file/vn-11134207-7qukw-ligbxxzzyv1r51", "Sold": "1,2k đã bán"},
			{"Name": "Giày Thể Thao Nam Nữ", "Price": "250.000đ", "Image": "https://down-vn.img.susercontent.com/file/vn-11134201-23030-x07l9p2x42ov7b", "Sold": "5k đã bán"},
			{"Name": "Balo Thời Trang Đi Học", "Price": "199.000đ", "Image": "https://down-vn.img.susercontent.com/file/vn-11134207-7qukw-lh1s6xeyv1r52c", "Sold": "3,4k đã bán"},
			{"Name": "Tai Nghe Bluetooth", "Price": "350.000đ", "Image": "https://down-vn.img.susercontent.com/file/vn-11134207-7r98o-lsth6k7xv1r52c", "Sold": "10k đã bán"},
			{"Name": "Đồng Hồ Thông Minh", "Price": "499.000đ", "Image": "https://down-vn.img.susercontent.com/file/vn-11134207-7r98o-lsth6k7y1r52c", "Sold": "2,1k đã bán"},
			{"Name": "Ốp Lưng Điện Thoại", "Price": "25.000đ", "Image": "https://down-vn.img.susercontent.com/file/vn-11134207-7r98o-lsth6k7z1r52c", "Sold": "50k đã bán"},
			{"Name": "Kem Chống Nắng", "Price": "120.000đ", "Image": "https://down-vn.img.susercontent.com/file/vn-11134207-7r98o-lsth6k7w1r52c", "Sold": "15k đã bán"},
			{"Name": "Mũ Lưỡi Trai Nam", "Price": "45.000đ", "Image": "https://down-vn.img.susercontent.com/file/vn-11134207-7r98o-lsth6k7v1r52c", "Sold": "8,2k đã bán"},
		},
	}

	err := h.tmpl.Execute(c.Writer, data)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
	}
}
