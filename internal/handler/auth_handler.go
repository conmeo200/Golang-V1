package handler

import (
	"encoding/json"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/auth"
	"github.com/conmeo200/Golang-V1/internal/dto"
	"github.com/conmeo200/Golang-V1/internal/logger"
	"github.com/conmeo200/Golang-V1/internal/service"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	service *service.AuthService
}

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.RespondWithError(w, dto.ErrInvalidRequest)
		logger.ErrorLogger.Printf("RegisterHandler: invalid request format")
		return
	}

	logger.AppLogger.Printf("Register Request received for %s", req.Email)

	user, err := h.service.RegisterUser(r.Context(), req.Email, req.Password)
	if err != nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, err.Error(), "REGISTRATION_FAILED"))
		return
	}

	// Generate Token here
	accessToken, refreshToken, err := auth.GenerateTokens(user.ID.String())
	if err != nil {
		dto.RespondWithError(w, dto.ErrInternal)
		return
	}

	dto.RespondWithSuccess(w, http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"user":          user,
		"refresh_token": refreshToken,
	}, "success")
}
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.RespondWithError(w, dto.ErrInvalidRequest)
		logger.ErrorLogger.Printf("LoginHandler: invalid request format")
		return
	}

	user, err := h.service.LoginUser(r.Context(), req.Email, req.Password)
	if err != nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusUnauthorized, err.Error(), "INVALID_CREDENTIALS"))
		logger.ErrorLogger.Printf("LoginHandler Unauthorized: %v", err)
		return
	}

	accessToken, refreshToken, err := auth.GenerateTokens(user.ID.String())
	if err != nil {
		dto.RespondWithError(w, dto.ErrInternal)
		return
	}

	dto.RespondWithSuccess(w, http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"user":          user,
		"refresh_token": refreshToken,
	}, "success")
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	// Typically we require the refresh token to revoke it in logout
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	// Decode may fail if body is empty, we ignore or handle
	json.NewDecoder(r.Body).Decode(&req)

	if req.RefreshToken != "" {
		// Attempt to revoke
		token, err := auth.ValidateToken(req.RefreshToken)
		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				exp, _ := claims["exp"].(float64)
				h.service.RevokeToken(r.Context(), req.RefreshToken, int64(exp))
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	dto.RespondWithSuccess(w, http.StatusOK, nil, "logged out successfuly")
}

func (h *AuthHandler) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.RespondWithError(w, dto.ErrInvalidRequest)
		logger.ErrorLogger.Printf("ForgotPasswordHandler decode error: %v", err)
		return
	}

	resetToken, err := h.service.ForgotPassword(r.Context(), req.Email)
	if err != nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, err.Error(), "FORGOT_PASSWORD_FAILED"))
		logger.ErrorLogger.Printf("ForgotPasswordHandler error: %v", err)
		return
	}

	logger.AppLogger.Printf("ForgotPassword token generated for %s", req.Email)

	dto.RespondWithSuccess(w, http.StatusOK, map[string]interface{}{
		"reset_token": resetToken,
	}, "reset token generated (mock email sent)")
}

func (h *AuthHandler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		dto.RespondWithError(w, dto.ErrUnauthorized)
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.RespondWithError(w, dto.ErrInvalidRequest)
		return
	}

	err := h.service.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, err.Error(), "PASSWORD_CHANGE_FAILED"))
		return
	}

	dto.RespondWithSuccess(w, http.StatusOK, nil, "password changed successfully")
}

func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.RespondWithError(w, dto.ErrInvalidRequest)
		logger.ErrorLogger.Printf("RefreshTokenHandler decode error: %v", err)
		return
	}

	accessToken, newRefreshToken, err := h.service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		dto.RespondWithError(w, dto.ErrUnauthorized)
		logger.ErrorLogger.Printf("RefreshTokenHandler error: %v", err)
		return
	}

	logger.AppLogger.Printf("Token refreshed successfully")

	logger.AppLogger.Printf("Token refreshed successfully")

	dto.RespondWithSuccess(w, http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	}, "token refreshed")
}

func (h *AuthHandler) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.RespondWithError(w, dto.ErrInvalidRequest)
		logger.ErrorLogger.Printf("RevokeTokenHandler decode error: %v", err)
		return
	}

	token, err := auth.ValidateToken(req.RefreshToken)
	if err == nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			exp, _ := claims["exp"].(float64)
			h.service.RevokeToken(r.Context(), req.RefreshToken, int64(exp))
			logger.AppLogger.Printf("Token revoked via API")
		} else {
			logger.ErrorLogger.Printf("RevokeTokenHandler claims format error")
		}
	} else {
		logger.ErrorLogger.Printf("RevokeTokenHandler validate error or invalid token: %v", err)
	}

	dto.RespondWithSuccess(w, http.StatusOK, nil, "token revoked successfully")
}