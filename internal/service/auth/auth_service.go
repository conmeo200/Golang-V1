package auth

import (
	"github.com/conmeo200/Golang-V1/internal/domain/model"
	repository "github.com/conmeo200/Golang-V1/internal/repository/auth"
	"github.com/google/uuid"
)
type AuthService struct {
	authRepo repository.AuthRepository
}
func NewAuthService(authRepository repository.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepository,
	}
}

func (a *AuthService)ProccessLogin() (model.User, error) {
	return a.authRepo.ProcessLogin()
}

func (a *AuthService)ProccessRegister() (model.User, error) {
	return a.authRepo.ProcessRegister()
	
}

func (a *AuthService)ProccessLogout() (uuid.UUID, error) {
	return a.authRepo.ProcessLogout()
}

func ProccessForgotPassWord(){
	
}

