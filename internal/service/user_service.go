package service

import (
	"context"
	"errors"

	"github.com/conmeo200/Golang-V1/internal/dto"
	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// userService is the concrete implementation of the UserService interface.
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new instance of the user service that implements UserService interface.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// CreateUser handles the logic for creating a new user.
func (s *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*model.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create the user model
	newUser := &model.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Balance:      req.Balance,
	}

	err = s.repo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// GetUserByID retrieves a user by their ID.
func (s *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetAllUsers retrieves all users.
func (s *userService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.GetAll(ctx)
}
