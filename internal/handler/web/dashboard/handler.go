package dashboard

import (
	"context"
	"net/http"
	"strings"

	"github.com/conmeo200/Golang-V1/internal/auth"
	"github.com/conmeo200/Golang-V1/internal/service"
	"github.com/golang-jwt/jwt/v5"
)

type DashboardHandler struct {
	authService service.AuthServiceInterface
	roleService service.RoleServiceInterface
}

func NewDashboardHandler(authService service.AuthServiceInterface, roleService service.RoleServiceInterface) *DashboardHandler {
	return &DashboardHandler{
		authService: authService,
		roleService: roleService,
	}
}

func (h *DashboardHandler) DashboardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenString := cookie.Value
		token, err := auth.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := h.authService.GetUserByID(r.Context(), userID)
		if err != nil || user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
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
