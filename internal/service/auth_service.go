package service

import (
	"context"
	"errors"
	"strings"

	"github.com/conmeo200/Golang-V1/internal/auth"
	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	RegisterUser(ctx context.Context, email string, password string) (*model.User, error)
	LoginUser(ctx context.Context, email string, password string) (*model.User, error)
	RevokeToken(ctx context.Context, tokenString string, expiresAt int64) error
	IsTokenBlacklisted(ctx context.Context, tokenString string) bool
	ChangePassword(ctx context.Context, userID string, oldPassword string, newPassword string) error
	ForgotPassword(ctx context.Context, email string) (string, error)
	RefreshToken(ctx context.Context, tokenString string) (string, string, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

type AuthService struct {
	authRepo  *repository.AuthRepository
	userRepo  repository.UserRepo
	tokenRepo repository.TokenRepo
}

func NewAuthService(
	authRepo *repository.AuthRepository,
	userRepo repository.UserRepo,
	tokenRepo repository.TokenRepo) *AuthService {
		return &AuthService{
			authRepo: authRepo,
			userRepo: userRepo,
			tokenRepo: tokenRepo,
		}
}

func (s *AuthService)RegisterUser(ctx context.Context, email string, password string) (*model.User, error) {
	email = strings.TrimSpace(email)

	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	// check existing user
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existing != nil {
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

	user, err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// return token
	return user, nil
}

func (s *AuthService) LoginUser(ctx context.Context, email string, password string) (*model.User, error) {
	email = strings.TrimSpace(email)

	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
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

func (s *AuthService) RevokeToken(ctx context.Context, tokenString string, expiresAt int64) error {
	blacklist := &model.TokenBlacklist{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}
	return s.tokenRepo.BlacklistToken(ctx, blacklist)
}

func (s *AuthService) IsTokenBlacklisted(ctx context.Context, tokenString string) bool {
	return s.tokenRepo.IsBlacklisted(ctx, tokenString)
}

func (s *AuthService) ChangePassword(ctx context.Context, userID string, oldPassword string, newPassword string) error {
	user, err := s.userRepo.GetUser(ctx, userID)
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

	return s.userRepo.UpdatePassword(ctx, userID, string(hash))
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) (string, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return "", errors.New("email is required")
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
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

func (s *AuthService) RefreshToken(ctx context.Context, tokenString string) (string, string, error) {
	if s.IsTokenBlacklisted(ctx, tokenString) {
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
	s.RevokeToken(ctx, tokenString, int64(exp))

	return auth.GenerateTokens(userID)
}

func (s *AuthService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.GetUser(ctx, id)
}