package service

import (
	"errors"
	"strings"

	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) FindFirstByEmail(email string) (*model.User, error) {

	email = strings.TrimSpace(email)

	if email == "" {
		return nil, errors.New("email is required")
	}

	user, err := s.repo.FindByEmail(email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateUser(email string, balance float64, password string) (*model.User, error) {

	email = strings.TrimSpace(email)

	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	if balance < 0 {
		return nil, errors.New("balance cannot be negative")
	}

	// check existing user
	existing, err := s.repo.FindByEmail(email)
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

	return s.repo.CreateUser(&user)
}

func (s *UserService) GetUser(id string) (*model.User, error) {

	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id is required")
	}

	return s.repo.GetUser(id)
}

func (s *UserService) ListUser() ([]model.User, error) {

	users, err := s.repo.ListUser()

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) UpdateBalance(id uint, newBalance float64) error {

	if id == 0 {
		return errors.New("invalid user id")
	}

	if newBalance < 0 {
		return errors.New("balance cannot be negative")
	}

	return s.repo.UpdateBalance(id, newBalance)
}

func (s *UserService) DeleteUser(id uint) error {

	if id == 0 {
		return errors.New("invalid user id")
	}

	return s.repo.Delete(id)
}