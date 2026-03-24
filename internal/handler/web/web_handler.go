package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/dto"
)

type WebHandler struct {
    // We can inject services here later if we want the web to talk to DB
}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

func (h *WebHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/login.html")
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (h *WebHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/register.html")
	if err != nil {
		log.Printf("Error parsing register template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

type DashboardData struct {
	Title      string
	ActiveMenu string
}

func (h *WebHandler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/layout.html", "web/template/dashboard.html")
	if err != nil {
		log.Printf("Error parsing dashboard template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := DashboardData{
		Title:      "Overview",
		ActiveMenu: "dashboard",
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		log.Printf("Error executing dashboard template: %v", err)
	}
}

type UserItem struct {
	ID       int
	Name     string
	Email    string
	Role     string
	Status   string
	JoinedAt string
}

type UserListPageData struct {
	Title      string
	ActiveMenu string
	Users      []UserItem
}

func (h *WebHandler) UserListPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/layout.html", "web/template/crud_list.html")
	if err != nil {
		log.Printf("Error parsing user list template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := UserListPageData{
		Title:      "Users Management",
		ActiveMenu: "users",
		Users: []UserItem{
			{ID: 1, Name: "Admin User", Email: "admin@heelo.com", Role: "Admin", Status: "Active", JoinedAt: "Jan 10, 2026"},
			{ID: 2, Name: "Editor Mike", Email: "mike@example.com", Role: "Editor", Status: "Active", JoinedAt: "Feb 15, 2026"},
			{ID: 3, Name: "Newbie Joe", Email: "joe@test.com", Role: "User", Status: "Pending", JoinedAt: "Mar 01, 2026"},
		},
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

type UserFormPageData struct {
	Title      string
	ActiveMenu string
	IsEdit     bool
	User       UserItem
}

func (h *WebHandler) UserFormPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/layout.html", "web/template/crud_form.html")
	if err != nil {
		log.Printf("Error parsing form template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := UserFormPageData{
		Title:      "Add User",
		ActiveMenu: "users",
		IsEdit:     false,
	}

	// Mocking edit mode if path contains 'edit'
	if r.URL.Path == "/dashboard/users/edit/1" {
		data.IsEdit = true
		data.Title = "Edit User"
		data.User = UserItem{ID: 1, Name: "Admin User", Email: "admin@heelo.com", Role: "Admin", Status: "Active"}
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *WebHandler) GradesPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/layout.html", "web/template/grade_list.html")
	if err != nil {
		log.Printf("Error parsing grades template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := DashboardData{
		Title:      "My Grades",
		ActiveMenu: "grades",
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *WebHandler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/layout.html", "web/template/student_profile.html")
	if err != nil {
		log.Printf("Error parsing profile template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Mock summary user data
	data := struct {
		Title      string
		ActiveMenu string
		User       struct {
			Name    string
			Initial string
			Email   string
		}
	}{
		Title:      "My Profile",
		ActiveMenu: "profile",
		User: struct {
			Name    string
			Initial string
			Email   string
		}{
			Name:    "Admin User",
			Initial: "AU",
			Email:   "admin@heelo.com",
		},
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *WebHandler) AnalyticsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/layout.html", "web/template/analytics.html")
	if err != nil {
		log.Printf("Error parsing analytics template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := DashboardData{
		Title:      "Academic Analytics",
		ActiveMenu: "analytics",
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

type NewsItem struct {
	Title    string
	Category string
	Excerpt  string
	Author   string
	Date     string
}

type NewsPageData struct {
	Title    string
	Heading  string
	NewsList []NewsItem
}

func (h *WebHandler) NewsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/news.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		dto.RespondWithError(w, dto.ErrInternal)
		return
	}

	data := NewsPageData{
		Title:   "HeeloNews - Latest Updates",
		Heading: "Latest Insights & Updates",
		NewsList: []NewsItem{
			{
				Title:    "Go 1.25 Released with Context Improvements",
				Category: "Technology",
				Excerpt:  "Latest Go version focus on concurrency and performance.",
				Author:   "Alice Wonderland",
				Date:     "Oct 15, 2026",
			},
		},
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func (h *WebHandler) SportsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/uniscore.html")
	if err != nil {
		log.Printf("Error parsing sports template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := DashboardData{
		Title:      "Live Sports Center",
		ActiveMenu: "sports",
	}

	// Because uniscore.html is a standalone page, we don't need layout.html
	tmpl.ExecuteTemplate(w, "uniscore.html", data)
}

func (h *WebHandler) MatchDetailPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/match_detail.html")
	if err != nil {
		log.Printf("Error parsing match detail template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := DashboardData{
		Title: "Match Detail - Real Madrid vs Atlético Madrid",
	}

	tmpl.ExecuteTemplate(w, "match_detail.html", data)
}

func (h *WebHandler) StandingsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/standings.html")
	if err != nil {
		log.Printf("Error parsing standings template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := DashboardData{
		Title: "Bảng xếp hạng - Uniscore",
	}

	tmpl.ExecuteTemplate(w, "standings.html", data)
}
