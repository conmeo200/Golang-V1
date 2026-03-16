package service

import (
	"context"
	"errors"

	"github.com/conmeo200/Golang-V1/internal/auth"
	"github.com/conmeo200/Golang-V1/internal/dto"
	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// authService is the concrete implementation of the AuthService interface.
type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository // Assuming TokenRepository interface exists
}

// NewAuthService creates a new instance that implements AuthService.
func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

// Register handles new user registration.
func (s *authService) Register(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if user already exists
	_, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		// If err is nil, it means a user was found
		return nil, errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Balance:      req.Balance, // Assuming balance comes from the request
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// Convert to DTO for the response
	userResponse := &dto.UserResponse{
		ID:      newUser.ID.String(),
		Email:   newUser.Email,
		Balance: newUser.Balance,
		Role:    newUser.Role,
		Status:  newUser.Status,
	}

	return userResponse, nil
}

// Login handles user authentication and token generation.
func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID.String()) // Simplified token generation
	if err != nil {
		return "", err
	}

	return token, nil
}

// NOTE: The following methods are kept from the old implementation.
// They should also be refactored to use interfaces and DTOs consistently.

func (s *authService) RevokeToken(ctx context.Context, tokenString string, expiresAt int64) error {
	blacklist := &model.TokenBlacklist{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}
	// This will cause an error if tokenRepo is not initialized with a concrete type
	// that has the BlacklistToken method.
	return s.tokenRepo.BlacklistToken(ctx, blacklist)
}

func (s *authService) IsTokenBlacklisted(ctx context.Context, tokenString string) bool {
	return s.tokenRepo.IsBlacklisted(ctx, tokenString)
}
