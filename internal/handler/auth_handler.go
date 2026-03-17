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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.SendError("invalid request format"))
		logger.ErrorLogger.Printf("RegisterHandler: invalid request format")
		return
	}

	logger.AppLogger.Printf("Register Request received for %s", req.Email)

	user, err := h.service.RegisterUser(r.Context(), req.Email, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.SendError(err.Error()))
		return
	}

	// Generate Token here
	accessToken, refreshToken, err := auth.GenerateTokens(user.ID.String())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.SendError("failed to generate token"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := dto.APIResponse{
		Status:  true,
		Message: "success",
		Data: map[string]interface{}{
			"access_token":  accessToken,
			"user":          user,
			"refresh_token": refreshToken,
		},
	}
	json.NewEncoder(w).Encode(resp)
}
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.SendError("invalid request format"))
		logger.ErrorLogger.Printf("LoginHandler: invalid request format")
		return
	}

	user, err := h.service.LoginUser(r.Context(), req.Email, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(dto.SendError(err.Error()))
		logger.ErrorLogger.Printf("LoginHandler Unauthorized: %v", err)
		return
	}

	accessToken, refreshToken, err := auth.GenerateTokens(user.ID.String())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.SendError("failed to generate token"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := dto.APIResponse{
		Status:  true,
		Message: "success",
		Data: map[string]interface{}{
			"access_token":  accessToken,
			"user":          user,
			"refresh_token": refreshToken,
		},
	}
	json.NewEncoder(w).Encode(resp)
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
	json.NewEncoder(w).Encode(dto.APIResponse{
		Status:  true,
		Message: "logged out successfuly",
	})
}

func (h *AuthHandler) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(dto.SendError("invalid request format"))
		logger.ErrorLogger.Printf("ForgotPasswordHandler decode error: %v", err)
		return
	}

	resetToken, err := h.service.ForgotPassword(r.Context(), req.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(dto.SendError(err.Error()))
		logger.ErrorLogger.Printf("ForgotPasswordHandler error: %v", err)
		return
	}

	logger.AppLogger.Printf("ForgotPassword token generated for %s", req.Email)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.APIResponse{
		Status:  true,
		Message: "reset token generated (mock email sent)",
		Data: map[string]interface{}{
			"reset_token": resetToken,
		},
	})
}

func (h *AuthHandler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(dto.SendError("unauthorized: user id not in context"))
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.SendError("invalid request format"))
		return
	}

	err := h.service.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.SendError(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.APIResponse{
		Status:  true,
		Message: "password changed successfully",
	})
}

func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(dto.SendError("invalid request format"))
		logger.ErrorLogger.Printf("RefreshTokenHandler decode error: %v", err)
		return
	}

	accessToken, newRefreshToken, err := h.service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(dto.SendError(err.Error()))
		logger.ErrorLogger.Printf("RefreshTokenHandler error: %v", err)
		return
	}

	logger.AppLogger.Printf("Token refreshed successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.APIResponse{
		Status:  true,
		Message: "token refreshed",
		Data: map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": newRefreshToken,
		},
	})
}

func (h *AuthHandler) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(dto.SendError("invalid request format"))
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.APIResponse{
		Status:  true,
		Message: "token revoked successfully",
	})
}