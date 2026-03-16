package service

import (
	"errors"
	"strings"

	"github.com/conmeo200/Golang-V1/internal/auth"
	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo *repository.AuthRepository
	userRepo *repository.UserRepository
	tokenRepo *repository.TokenRepository
}

func NewAuthService(
	authRepo *repository.AuthRepository,
	userRepo *repository.UserRepository,
	tokenRepo *repository.TokenRepository) *AuthService {
		return &AuthService{
			authRepo: authRepo,
			userRepo: userRepo,
			tokenRepo: tokenRepo,
		}
}

func (s *AuthService)RegisterUser(email string, password string) (*model.User, error) {
	email = strings.TrimSpace(email)

	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	// check existing user
	existing, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("email already exists")
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email		 :    email,
		PasswordHash : string(hash),
	}

	user, err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// return token
	return user, nil
}

func (s *AuthService) LoginUser(email string, password string) (*model.User, error) {
	email = strings.TrimSpace(email)

	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !CheckPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func CheckPassword(password string, password_hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	return err == nil
}

func (s *AuthService) RevokeToken(tokenString string, expiresAt int64) error {
	blacklist := &model.TokenBlacklist{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}
	return s.tokenRepo.BlacklistToken(blacklist)
}

func (s *AuthService) IsTokenBlacklisted(tokenString string) bool {
	return s.tokenRepo.IsBlacklisted(tokenString)
}

func (s *AuthService) ChangePassword(userID string, oldPassword string, newPassword string) error {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if !CheckPassword(oldPassword, user.PasswordHash) {
		return errors.New("invalid old password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, string(hash))
}

func (s *AuthService) ForgotPassword(email string) (string, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return "", errors.New("email is required")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	// Mocking email sending by returning a reset token
	resetToken := "mock_reset_token_" + user.ID.String()
	return resetToken, nil
}

func (s *AuthService) RefreshToken(tokenString string) (string, string, error) {
	if s.IsTokenBlacklisted(tokenString) {
		return "", "", errors.New("token is blacklisted")
	}

	token, err := auth.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("invalid user id in token")
	}

	// Revoke old refresh token so it can't be used again
	exp, _ := claims["exp"].(float64)
	s.RevokeToken(tokenString, int64(exp))

	return auth.GenerateTokens(userID)
}