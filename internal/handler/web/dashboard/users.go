package dashboard

import (
	"html/template"
	"log"
	"net/http"
)

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

func (h *DashboardHandler) UserListPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/crud_list.html")
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

func (h *DashboardHandler) UserFormPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/crud_form.html")
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

	if r.URL.Path == "/dashboard/users/edit/1" {
		data.IsEdit = true
		data.Title = "Edit User"
		data.User = UserItem{ID: 1, Name: "Admin User", Email: "admin@heelo.com", Role: "Admin", Status: "Active"}
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *DashboardHandler) UserDetailPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/user_detail.html")
	if err != nil {
		log.Printf("Error parsing user detail template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := UserFormPageData{
		Title:      "User Details",
		ActiveMenu: "users",
		User:       UserItem{ID: 1, Name: "Admin User", Email: "admin@heelo.com", Role: "Admin", Status: "Active", JoinedAt: "Jan 10, 2026"},
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}
