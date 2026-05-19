package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type RoleItem struct {
	ID          int
	Name        string
	Description string
	UsersCount  int
}

type RoleListPageData struct {
	Title      string
	ActiveMenu string
	Roles      []RoleItem
}

func (h *DashboardHandler) RoleListPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/role_list.html")
	if err != nil {
		log.Printf("Error parsing role list template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	roles, err := h.roleService.GetAllRolesWithUserCount()
	if err != nil {
		http.Error(w, "Error fetching roles", http.StatusInternalServerError)
		return
	}

	var roleItems []RoleItem
	for _, role := range roles {
		roleItems = append(roleItems, RoleItem{
			ID:          int(role.ID),
			Name:        role.Name,
			Description: role.Description,
			UsersCount:  role.UsersCount,
		})
	}

	data := RoleListPageData{
		Title:      "Role Management",
		ActiveMenu: "roles",
		Roles:      roleItems,
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

type PermissionItem struct {
	ID          string
	Action      string
	Description string
	Granted     bool
}

type PermissionGroup struct {
	Module      string
	Description string
	Permissions []PermissionItem
}

type RoleDetailPageData struct {
	Title            string
	ActiveMenu       string
	Role             RoleItem
	PermissionGroups []PermissionGroup
	Message          string
}

func (h *DashboardHandler) RoleDetailPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/role_detail.html")
	if err != nil {
		log.Printf("Error parsing role detail template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	roleIDStr := pathParts[len(pathParts)-1]
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid Role ID", http.StatusBadRequest)
		return
	}

	message := ""
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err == nil {
			selectedPerms := r.Form["permissions"]
			h.roleService.UpdateRolePermissions(uint(roleID), selectedPerms)
			message = "Permissions updated successfully!"
		}
	}

	role, err := h.roleService.GetRoleWithPermissions(uint(roleID))
	if err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	allPerms, err := h.roleService.GetAllPermissions()
	if err != nil {
		http.Error(w, "Error fetching permissions", http.StatusInternalServerError)
		return
	}

	grantedPermMap := make(map[string]bool)
	for _, p := range role.Permissions {
		grantedPermMap[p.ID] = true
	}

	modulePermMap := make(map[string][]PermissionItem)
	for _, p := range allPerms {
		modulePermMap[p.Module] = append(modulePermMap[p.Module], PermissionItem{
			ID:          p.ID,
			Action:      p.Action,
			Description: p.Description,
			Granted:     grantedPermMap[p.ID],
		})
	}

	var groups []PermissionGroup
	for module, perms := range modulePermMap {
		groups = append(groups, PermissionGroup{
			Module:      module,
			Description: "Manage " + module,
			Permissions: perms,
		})
	}

	data := RoleDetailPageData{
		Title:            "Role Permissions - " + role.Name,
		ActiveMenu:       "roles",
		Role:             RoleItem{ID: int(role.ID), Name: role.Name, Description: role.Description},
		PermissionGroups: groups,
		Message:          message,
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

type RoleFormPageData struct {
	Title      string
	ActiveMenu string
	IsEdit     bool
	Role       RoleItem
}

func (h *DashboardHandler) RoleFormPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/role_form.html")
	if err != nil {
		log.Printf("Error parsing role form template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := RoleFormPageData{
		Title:      "Create New Role",
		ActiveMenu: "roles",
		IsEdit:     false,
	}

	// Check if edit mode
	if strings.Contains(r.URL.Path, "/edit/") {
		pathParts := strings.Split(r.URL.Path, "/")
		roleIDStr := pathParts[len(pathParts)-1]
		roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
		if err == nil {
			role, err := h.roleService.GetRoleWithPermissions(uint(roleID))
			if err == nil {
				data.IsEdit = true
				data.Title = "Edit Role"
				data.Role = RoleItem{
					ID:          int(role.ID),
					Name:        role.Name,
					Description: role.Description,
				}
			}
		}
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *DashboardHandler) ProcessRoleCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	err = h.roleService.CreateRole(name, description)
	if err != nil {
		http.Error(w, "Failed to create role", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard/roles", http.StatusSeeOther)
}

func (h *DashboardHandler) ProcessRoleUpdate(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	roleIDStr := pathParts[len(pathParts)-1]
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid Role ID", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	err = h.roleService.UpdateRole(uint(roleID), name, description)
	if err != nil {
		http.Error(w, "Failed to update role", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard/roles", http.StatusSeeOther)
}

func (h *DashboardHandler) ProcessRoleDelete(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	roleIDStr := pathParts[len(pathParts)-1]
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid Role ID", http.StatusBadRequest)
		return
	}

	err = h.roleService.DeleteRole(uint(roleID))
	if err != nil {
		http.Error(w, "Failed to delete role", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard/roles", http.StatusSeeOther)
}
