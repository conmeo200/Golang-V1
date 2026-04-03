package dashboard

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/google/uuid"
)

// TaxListPageData holds data for the tax list page
type TaxListPageData struct {
	Title      string
	ActiveMenu string
	Taxes      interface{} // We'll use a specific slice type later
	Error      string
}

// TaxListPage renders the tax declarations list
func (h *DashboardHandler) TaxListPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/tax_list.html")
	if err != nil {
		log.Printf("Error parsing tax list template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// For now, use a placeholder userID. In production, get from context.
	userID, _ := uuid.Parse("00000000-0000-0000-0000-000000000001") // Placeholder
	
	decls, err := h.taxService.ListDeclarations(r.Context(), userID)
	if err != nil {
		log.Printf("Error listing declarations: %v", err)
	}

	data := TaxListPageData{
		Title:      "Tax Declarations",
		ActiveMenu: "taxes",
		Taxes:      decls,
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

// TaxNewPage renders the form to create a new tax declaration
func (h *DashboardHandler) TaxNewPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/tax_form.html")
	if err != nil {
		log.Printf("Error parsing tax form template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := TaxListPageData{
		Title:      "New Tax Declaration",
		ActiveMenu: "taxes",
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

// ProcessTaxDeclaration handles the form submission for a new tax declaration
func (h *DashboardHandler) ProcessTaxDeclaration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Placeholder UserID
	userID, _ := uuid.Parse("00000000-0000-0000-0000-000000000001") // Placeholder

	var totalIncome int64
	fmt.Sscanf(r.FormValue("income_amount"), "%d", &totalIncome)

	decl := &model.TaxDeclaration{
		UserID:      userID,
		Type:        model.TaxType(r.FormValue("type")),
		Period:      r.FormValue("period"),
		TotalIncome: totalIncome,
		Status:      model.StatusSubmitted,
	}

	err = h.taxService.CreateDeclaration(r.Context(), decl)
	if err != nil {
		log.Printf("Error creating declaration: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard/taxes", http.StatusSeeOther)
}

