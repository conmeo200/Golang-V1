package service

import (
	"context"
	"errors"
	"strings"

	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	FindFirstByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, email string, balance float64, password string) (*model.User, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
	ListUser(ctx context.Context) ([]model.User, error)
	UpdateBalance(ctx context.Context, id uint, newBalance float64) error
	DeleteUser(ctx context.Context, id uint) error
}

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) FindFirstByEmail(ctx context.Context, email string) (*model.User, error) {

	email = strings.TrimSpace(email)

	if email == "" {
		return nil, errors.New("email is required")
	}

	user, err := s.repo.FindByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, email string, balance float64, password string) (*model.User, error) {

	email = strings.TrimSpace(email)

	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	// check existing user
	existing, err := s.repo.FindByEmail(ctx, email)
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

	user := model.User{
		Email:        email,
		Balance:      balance,
		PasswordHash: string(hash),
	}

	return s.repo.CreateUser(ctx, &user)
}

func (s *UserService) GetUser(ctx context.Context, id string) (*model.User, error) {

	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id is required")
	}

	return s.repo.GetUser(ctx, id)
}

func (s *UserService) ListUser(ctx context.Context) ([]model.User, error) {

	users, err := s.repo.ListUser(ctx)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) UpdateBalance(ctx context.Context, id uint, newBalance float64) error {

	if id == 0 {
		return errors.New("invalid user id")
	}

	if newBalance < 0 {
		return errors.New("balance cannot be negative")
	}

	return s.repo.UpdateBalance(ctx, id, newBalance)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {

	if id == 0 {
		return errors.New("invalid user id")
	}

	return s.repo.Delete(ctx, id)
}