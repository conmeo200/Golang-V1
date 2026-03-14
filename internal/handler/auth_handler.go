package handler

import (
	//"encoding/json"
	"encoding/json"
	"log"
	"net/http"

	//"github.com/conmeo200/Golang-V1/internal/auth"
	//"github.com/conmeo200/Golang-V1/internal/dto"
	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

type Request struct {
	Email string 
	Password  string
}

func CheckPassword(password string, user *model.User) bool {
	return true
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler Register")
	var req Request
	json.NewDecoder(r.Body).Decode(&req)

	log.Println("Request", req)
	// user, err := h.serviceUser.FindByEmail(req.Email)
	// if err != nil {
	// 	http.Error(w, "user not found", 401)
	// 	return
	// }

	// if !CheckPassword(req.Password, user.PasswordHash) {
	// 	http.Error(w, "invalid password", 401)
	// 	return
	// }

	// token, _ := auth.GenerateToken(user.ID.String())

	// resp := dto.APIResponse{
	// 	Status: true,
	// 	Data: map[string]string{
	// 		"token": token,
	// 	},
	// }

	// json.NewEncoder(w).Encode(resp)
}
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	// var req LoginRequest
	// json.NewDecoder(r.Body).Decode(&req)

	// user, err := h.serviceUser.FindByEmail(req.Email)
	// if err != nil {
	// 	http.Error(w, "user not found", 401)
	// 	return
	// }

	// if !CheckPassword(req.Password, user.PasswordHash) {
	// 	http.Error(w, "invalid password", 401)
	// 	return
	// }

	// token, _ := auth.GenerateToken(user.ID.String())

	// resp := dto.APIResponse{
	// 	Status: true,
	// 	Data: map[string]string{
	// 		"token": token,
	// 	},
	// }

	// json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	// var req LoginRequest
	// json.NewDecoder(r.Body).Decode(&req)

	// user, err := h.serviceUser.FindByEmail(req.Email)
	// if err != nil {
	// 	http.Error(w, "user not found", 401)
	// 	return
	// }

	// if !CheckPassword(req.Password, user.PasswordHash) {
	// 	http.Error(w, "invalid password", 401)
	// 	return
	// }

	// token, _ := auth.GenerateToken(user.ID.String())

	// resp := dto.APIResponse{
	// 	Status: true,
	// 	Data: map[string]string{
	// 		"token": token,
	// 	},
	// }

	// json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) ForgetPasswordHandler(w http.ResponseWriter, r *http.Request) {

	// var req LoginRequest
	// json.NewDecoder(r.Body).Decode(&req)

	// user, err := h.serviceUser.FindByEmail(req.Email)
	// if err != nil {
	// 	http.Error(w, "user not found", 401)
	// 	return
	// }

	// if !CheckPassword(req.Password, user.PasswordHash) {
	// 	http.Error(w, "invalid password", 401)
	// 	return
	// }

	// token, _ := auth.GenerateToken(user.ID.String())

	// resp := dto.APIResponse{
	// 	Status: true,
	// 	Data: map[string]string{
	// 		"token": token,
	// 	},
	// }

	// json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {

	// var req LoginRequest
	// json.NewDecoder(r.Body).Decode(&req)

	// user, err := h.serviceUser.FindByEmail(req.Email)
	// if err != nil {
	// 	http.Error(w, "user not found", 401)
	// 	return
	// }

	// if !CheckPassword(req.Password, user.PasswordHash) {
	// 	http.Error(w, "invalid password", 401)
	// 	return
	// }

	// token, _ := auth.GenerateToken(user.ID.String())

	// resp := dto.APIResponse{
	// 	Status: true,
	// 	Data: map[string]string{
	// 		"token": token,
	// 	},
	// }

	// json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	// var req LoginRequest
	// json.NewDecoder(r.Body).Decode(&req)

	// user, err := h.serviceUser.FindByEmail(req.Email)
	// if err != nil {
	// 	http.Error(w, "user not found", 401)
	// 	return
	// }

	// if !CheckPassword(req.Password, user.PasswordHash) {
	// 	http.Error(w, "invalid password", 401)
	// 	return
	// }

	// token, _ := auth.GenerateToken(user.ID.String())

	// resp := dto.APIResponse{
	// 	Status: true,
	// 	Data: map[string]string{
	// 		"token": token,
	// 	},
	// }

	// json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {

	// var req LoginRequest
	// json.NewDecoder(r.Body).Decode(&req)

	// user, err := h.serviceUser.FindByEmail(req.Email)
	// if err != nil {
	// 	http.Error(w, "user not found", 401)
	// 	return
	// }

	// if !CheckPassword(req.Password, user.PasswordHash) {
	// 	http.Error(w, "invalid password", 401)
	// 	return
	// }

	// token, _ := auth.GenerateToken(user.ID.String())

	// resp := dto.APIResponse{
	// 	Status: true,
	// 	Data: map[string]string{
	// 		"token": token,
	// 	},
	// }

	// json.NewEncoder(w).Encode(resp)
}