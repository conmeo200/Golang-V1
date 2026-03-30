package dashboard

import (
	"html/template"
	"log"
	"net/http"
)

type DashboardData struct {
	Title      string
	ActiveMenu string
}

func (h *DashboardHandler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/dashboard.html")
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

func (h *DashboardHandler) GradesPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/grade_list.html")
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

func (h *DashboardHandler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/student_profile.html")
	if err != nil {
		log.Printf("Error parsing profile template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

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

func (h *DashboardHandler) AnalyticsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/analytics.html")
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
