package dashboard

import (
	"html/template"
	"log"
	"net/http"
)

func (h *DashboardHandler) DashboardLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/login.html")
	if err != nil {
		log.Printf("Error parsing dashboard login template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]interface{}{
		"Error": r.URL.Query().Get("error"),
	})
}

func (h *DashboardHandler) ProcessDashboardLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email 	 := r.FormValue("email")
	password := r.FormValue("password")

	_, err := h.authService.LoginUser(r.Context(), email, password)
	if err != nil {
		http.Redirect(w, r, "/dashboard/login?error=invalid_credentials", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
