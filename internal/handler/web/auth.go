package web

import (
	"html/template"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/auth"
	"github.com/conmeo200/Golang-V1/internal/module/auth/port"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService port.AuthService
}

func NewAuthHandler(service port.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: service,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	data := map[string]interface{}{
		"Title": "Đăng nhập",
		"Error": c.Query("error"),
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
		"Error": c.Query("error"),
	}
	tmpl, err := template.ParseFiles("web/templates/pages/auth/register.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.Execute(c.Writer, data)
}

func (h *AuthHandler) LoginPost(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := h.authService.LoginUser(c.Request.Context(), email, password)
	if err != nil {
		c.Redirect(http.StatusFound, "/login?error=invalid_credentials")
		return
	}

	accessToken, _, err := auth.GenerateTokens(user.ID.String())
	if err != nil {
		c.Redirect(http.StatusFound, "/login?error=server_error")
		return
	}

	// Just use a simple cookie for the web session
	c.SetCookie("session_token", accessToken, 3600*24, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func (h *AuthHandler) RegisterPost(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := h.authService.RegisterUser(c.Request.Context(), email, password)
	if err != nil {
		c.Redirect(http.StatusFound, "/register?error=registration_failed")
		return
	}

	accessToken, _, err := auth.GenerateTokens(user.ID.String())
	if err != nil {
		c.Redirect(http.StatusFound, "/register?error=server_error")
		return
	}

	c.SetCookie("session_token", accessToken, 3600*24, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}


