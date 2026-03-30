package service

import (
	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/repository"
)

type RoleServiceInterface interface {
	GetAllRolesWithUserCount() ([]RoleWithUserCount, error)
	GetRoleWithPermissions(id uint) (*model.Role, error)
	GetAllPermissions() ([]model.Permission, error)
	UpdateRolePermissions(roleID uint, permissionIDs []string) error
}

type RoleService struct {
	repo repository.RoleRepositoryInterface
}

func NewRoleService(repo repository.RoleRepositoryInterface) RoleServiceInterface {
	return &RoleService{repo: repo}
}

type RoleWithUserCount struct {
	model.Role
	UsersCount int
}

func (s *RoleService) GetAllRolesWithUserCount() ([]RoleWithUserCount, error) {
	roles, err := s.repo.GetAllRoles()
	if err != nil {
		return nil, err
	}

	var result []RoleWithUserCount
	for _, r := range roles {
		count, _ := s.repo.GetUsersCountByRoleName(r.Name)
		result = append(result, RoleWithUserCount{
			Role:       r,
			UsersCount: int(count),
		})
	}

	return result, nil
}

func (s *RoleService) GetRoleWithPermissions(id uint) (*model.Role, error) {
	return s.repo.GetRoleWithPermissions(id)
}

func (s *RoleService) GetAllPermissions() ([]model.Permission, error) {
	return s.repo.GetAllPermissions()
}

func (s *RoleService) UpdateRolePermissions(roleID uint, permissionIDs []string) error {
	return s.repo.UpdateRolePermissions(roleID, permissionIDs)
}
