package service

import "github.com/conmeo200/Golang-V1/internal/repository"

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repo: r}
}