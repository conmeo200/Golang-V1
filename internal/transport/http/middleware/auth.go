package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/conmeo200/Golang-V1/internal/auth"
	authmodule "github.com/conmeo200/Golang-V1/internal/module/auth"
	"github.com/conmeo200/Golang-V1/internal/module/role"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	authService authmodule.AuthServiceInterface
	roleService role.RoleServiceInterface
}

func NewAuthMiddleware(authService authmodule.AuthServiceInterface, roleService role.RoleServiceInterface) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		roleService: roleService,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Redirect(w, r, "/dashboard/login", http.StatusSeeOther)
			return
		}

		tokenString := cookie.Value
		token, err := auth.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			http.Redirect(w, r, "/dashboard/login", http.StatusSeeOther)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Redirect(w, r, "/dashboard/login", http.StatusSeeOther)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Redirect(w, r, "/dashboard/login", http.StatusSeeOther)
			return
		}

		user, err := m.authService.GetUserByID(r.Context(), userID)
		if err != nil || user == nil {
			http.Redirect(w, r, "/dashboard/login", http.StatusSeeOther)
			return
		}

		if strings.ToLower(user.Role) != "admin" {
			http.Error(w, "Forbidden - Admin access required", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequirePermission(permissionID string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			http.Error(w, "Unauthorized - Missing User Context", http.StatusUnauthorized)
			return
		}

		user, err := m.authService.GetUserByID(r.Context(), userID)
		if err != nil || user == nil {
			http.Error(w, "Unauthorized - User Not Found", http.StatusUnauthorized)
			return
		}

		// Admin string override for backward compatibility
		if strings.ToLower(user.Role) == "admin" {
			next.ServeHTTP(w, r)
			return
		}

		if user.RoleID == 0 {
			http.Error(w, "Forbidden - No Role Assigned", http.StatusForbidden)
			return
		}

		roleEntity, err := m.roleService.GetRoleWithPermissions(user.RoleID)
		if err != nil || roleEntity == nil {
			http.Error(w, "Forbidden - Role Not Found", http.StatusForbidden)
			return
		}

		hasPerm := false
		for _, p := range roleEntity.Permissions {
			if p.ID == permissionID {
				hasPerm = true
				break
			}
		}

		if !hasPerm {
			http.Error(w, "Forbidden - Missing Permission: "+permissionID, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
