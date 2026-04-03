package dashboard

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/conmeo200/Golang-V1/internal/service"
)

// LogListPageData holds data for the log list page
type LogListPageData struct {
	Title      string
	ActiveMenu string
	Logs       []service.LogFileInfo
	Total      int
	Error      string
}

// LogListPage renders the log file list page
func (h *DashboardHandler) LogListPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/log_list.html")
	if err != nil {
		log.Printf("Error parsing log list template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logs, err := h.logService.ListLogs()
	if err != nil {
		tmpl.ExecuteTemplate(w, "layout.html", LogListPageData{
			Title:      "System Logs",
			ActiveMenu: "logs",
			Error:      "Could not retrieve logs: " + err.Error(),
		})
		return
	}

	data := LogListPageData{
		Title:      "System Logs",
		ActiveMenu: "logs",
		Logs:       logs,
		Total:      len(logs),
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

// LogDetailPageData holds data for the log detail page
type LogDetailPageData struct {
	Title      string
	ActiveMenu string
	Filename   string
	Content    template.HTML
	Error      string
}

// LogDetailPage renders the log file content page
func (h *DashboardHandler) LogDetailPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/log_detail.html")
	if err != nil {
		log.Printf("Error parsing log detail template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Extract the filename from the URL, expecting `/dashboard/logs/detail/{filename}`
	pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathSegments) < 4 {
		http.Redirect(w, r, "/dashboard/logs", http.StatusSeeOther)
		return
	}

	filename := pathSegments[3]

	content, err := h.logService.GetLogContent(filename)
	if err != nil {
		tmpl.ExecuteTemplate(w, "layout.html", LogDetailPageData{
			Title:      "Log Content",
			ActiveMenu: "logs",
			Filename:   filename,
			Error:      "Could not read log file: " + err.Error(),
		})
		return
	}

	data := LogDetailPageData{
		Title:      "Log Content - " + filename,
		ActiveMenu: "logs",
		Filename:   filename,
		Content:    template.HTML(strings.ReplaceAll(template.HTMLEscapeString(content), "\n", "<br>")),
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

