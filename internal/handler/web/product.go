package web

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (h *ProductHandler) Detail(c *gin.Context) {
	// In a real app, we would fetch product by c.Param("id")
	data := map[string]interface{}{
		"Title": "Chi tiết sản phẩm",
		"Product": map[string]interface{}{
			"ID":          1,
			"Name":        "Áo Thun Nam Cổ Tròn Premium Cotton Chống Nhăn Tuyệt Đối",
			"Price":       "150.000đ",
			"Category":    "Thời trang Nam",
			"Image":       "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image",
			"Sold":        "1,2k",
			"Reviews":     456,
			"Description": "Áo thun nam cổ tròn được dệt từ 100% sợi cotton tự nhiên mang lại cảm giác mềm mại, thoáng mát và co giãn tốt. Thiết kế form dáng basic phù hợp với mọi vóc dáng, dễ dàng phối hợp cùng nhiều phong cách khác nhau từ trẻ trung, năng động đến thanh lịch. Đặc biệt công nghệ chống nhăn giúp áo luôn giữ được form dáng hoàn hảo sau nhiều lần giặt.",
		},
	}

	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/pages/product/detail.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *ProductHandler) Category(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Thời trang Nam",
		"Products": []map[string]interface{}{
			{"ID": 1, "Name": "Áo Thun Nam Cổ Tròn Premium Cotton", "Price": "150.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Sold": "1,2k đã bán"},
			{"ID": 2, "Name": "Giày Thể Thao Nam Nữ Siêu Nhẹ", "Price": "250.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Sold": "5k đã bán"},
			{"ID": 8, "Name": "Mũ Lưỡi Trai Nam Nữ Basic", "Price": "45.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Sold": "8,2k đã bán"},
			{"ID": 1, "Name": "Áo Thun Nam Cổ Tròn Premium Cotton", "Price": "150.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Sold": "1,2k đã bán"},
			{"ID": 2, "Name": "Giày Thể Thao Nam Nữ Siêu Nhẹ", "Price": "250.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Sold": "5k đã bán"},
			{"ID": 8, "Name": "Mũ Lưỡi Trai Nam Nữ Basic", "Price": "45.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "Sold": "8,2k đã bán"},
		},
	}

	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/pages/product/category.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}
