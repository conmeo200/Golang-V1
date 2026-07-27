package web

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type AffiliateHandler struct{}

func NewAffiliateHandler() *AffiliateHandler {
	return &AffiliateHandler{}
}

func (h *AffiliateHandler) Dashboard(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Affiliate Dashboard",
		"ActiveMenu": "dashboard",
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/affiliate_base.html", "web/templates/pages/affiliate/dashboard.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template: " + err.Error())
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "affiliate_base", data)
}

func (h *AffiliateHandler) Links(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Tracking Links",
		"ActiveMenu": "links",
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/affiliate_base.html", "web/templates/pages/affiliate/links.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template: " + err.Error())
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "affiliate_base", data)
}

func (h *AffiliateHandler) Performance(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Hiệu suất",
		"ActiveMenu": "performance",
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/affiliate_base.html", "web/templates/pages/affiliate/performance.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template: " + err.Error())
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "affiliate_base", data)
}

func (h *AffiliateHandler) Orders(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Báo cáo Đơn hàng",
		"ActiveMenu": "orders",
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/affiliate_base.html", "web/templates/pages/affiliate/orders.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template: " + err.Error())
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "affiliate_base", data)
}

func (h *AffiliateHandler) Payment(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Thanh toán",
		"ActiveMenu": "payment",
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/affiliate_base.html", "web/templates/pages/affiliate/payment.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template: " + err.Error())
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "affiliate_base", data)
}
