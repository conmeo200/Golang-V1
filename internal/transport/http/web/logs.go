package web

import (
	"net/http"
)

func (h *DashboardHandler) LogListPage(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Logs functionality temporarily disabled during refactor", http.StatusNotImplemented)
}

func (h *DashboardHandler) LogDetailPage(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Logs functionality temporarily disabled during refactor", http.StatusNotImplemented)
}
