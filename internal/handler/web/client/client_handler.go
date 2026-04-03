package client

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/conmeo200/Golang-V1/internal/auth"
	"github.com/conmeo200/Golang-V1/internal/dto"
	"github.com/conmeo200/Golang-V1/internal/service"
)

type ClientHandler struct {
	authService service.AuthServiceInterface
}

func NewClientHandler(authService service.AuthServiceInterface) *ClientHandler {
	return &ClientHandler{
		authService: authService,
	}
}

func (h *ClientHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/client/login.html")
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]interface{}{
		"Error": r.URL.Query().Get("error"),
	})
}

func (h *ClientHandler) ProcessClientLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.authService.LoginUser(r.Context(), email, password)
	if err != nil {
		http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
		return
	}

	accessToken, _, err := auth.GenerateTokens(user.ID.String())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(15 * time.Minute),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *ClientHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *ClientHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/client/register.html")
	if err != nil {
		log.Printf("Error parsing register template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

type HomePageData struct {
	Title   string
	Heading string
}

func (h *ClientHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/client/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		dto.RespondWithError(w, dto.ErrInternal)
		return
	}

	data := HomePageData{
		Title:   "HeeloApp - Welcome",
		Heading: "Welcome to HeeloApp",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

type PageData struct {
	Title      string
	ActiveMenu string
}

func (h *ClientHandler) SportsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/client/uniscore.html")
	if err != nil {
		log.Printf("Error parsing sports template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:      "Live Sports Center",
		ActiveMenu: "sports",
	}
	tmpl.ExecuteTemplate(w, "uniscore.html", data)
}

func (h *ClientHandler) MatchDetailPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/client/match_detail.html")
	if err != nil {
		log.Printf("Error parsing match detail template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title: "Match Detail - Real Madrid vs Atlético Madrid",
	}
	tmpl.ExecuteTemplate(w, "match_detail.html", data)
}

func (h *ClientHandler) StandingsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/client/standings.html")
	if err != nil {
		log.Printf("Error parsing standings template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title: "Bảng xếp hạng - Uniscore",
	}
	tmpl.ExecuteTemplate(w, "standings.html", data)
}
