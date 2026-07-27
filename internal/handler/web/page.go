package web

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type PageHandler struct{}

func NewPageHandler() *PageHandler {
	return &PageHandler{}
}

func (h *PageHandler) About(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Giới thiệu",
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/pages/info/about.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *PageHandler) Contact(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Liên hệ",
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/pages/info/contact.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *PageHandler) Blog(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Tin tức thời trang",
		"Posts": []map[string]interface{}{
			{"Title": "Xu hướng thời trang Thu Đông 2026", "Image": "https://placehold.co/600x400/E2E8F0/64748B?text=Fashion+1", "Date": "01/07/2026", "Excerpt": "Khám phá những bộ sưu tập mới nhất với phong cách tối giản nhưng vô cùng sang trọng."},
			{"Title": "Cách phối đồ với áo khoác Blazer", "Image": "https://placehold.co/600x400/E2E8F0/64748B?text=Fashion+2", "Date": "28/06/2026", "Excerpt": "Blazer không bao giờ lỗi mốt, nhưng làm sao để mặc đẹp nhất? Hãy cùng tìm hiểu."},
			{"Title": "Top 5 màu sắc dẫn đầu xu hướng năm nay", "Image": "https://placehold.co/600x400/E2E8F0/64748B?text=Fashion+3", "Date": "15/06/2026", "Excerpt": "Màu sắc quyết định cá tính của bạn. Dưới đây là những màu sắc nổi bật nhất."},
		},
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/pages/info/blog.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *PageHandler) Promotions(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Khuyến mãi",
		"Vouchers": []map[string]interface{}{
			{"Code": "SUMMER26", "Discount": "Giảm 20%", "Desc": "Áp dụng cho đơn từ 500k", "ValidUntil": "31/07/2026"},
			{"Code": "FREESHIP", "Discount": "Freeship", "Desc": "Áp dụng cho mọi đơn hàng", "ValidUntil": "15/07/2026"},
			{"Code": "NEWUSER", "Discount": "Giảm 50k", "Desc": "Khách hàng mới", "ValidUntil": "Không giới hạn"},
		},
		"Products": []map[string]interface{}{
			{"ID": 10, "Name": "Áo Khoác Nam Chống Nước", "Price": "300.000đ", "OldPrice": "500.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Sale+1", "Sold": "1,2k đã bán"},
			{"ID": 11, "Name": "Túi Xách Nữ Da Thật", "Price": "450.000đ", "OldPrice": "900.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Sale+2", "Sold": "5k đã bán"},
			{"ID": 12, "Name": "Set Đồ Tập Gym Nữ", "Price": "199.000đ", "OldPrice": "400.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Sale+3", "Sold": "3,4k đã bán"},
			{"ID": 13, "Name": "Tai Nghe Bluetooth Pro", "Price": "350.000đ", "OldPrice": "700.000đ", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Sale+4", "Sold": "10k đã bán"},
		},
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/pages/info/promotions.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *PageHandler) NotFound(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "404 - Không tìm thấy trang",
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/pages/error/404.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	c.Status(http.StatusNotFound)
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}
