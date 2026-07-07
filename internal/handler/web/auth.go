package web

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Đăng nhập",
	}
	tmpl, err := template.ParseFiles("web/templates/pages/auth/login.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.Execute(c.Writer, data)
}

func (h *AuthHandler) Register(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Đăng ký",
	}
	tmpl, err := template.ParseFiles("web/templates/pages/auth/register.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.Execute(c.Writer, data)
}
