package role

import (
	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/persistence"
)

type RoleServiceInterface interface {
	GetAllRolesWithUserCount() ([]RoleWithUserCount, error)
	GetRoleWithPermissions(id uint) (*model.Role, error)
	GetAllPermissions() ([]model.Permission, error)
	UpdateRolePermissions(roleID uint, permissionIDs []string) error
	CreateRole(name, description string) error
	UpdateRole(id uint, name, description string) error
	DeleteRole(id uint) error
	SeedDefaultPermissions() error
}

type RoleService struct {
	repo persistence.RoleRepositoryInterface
}

func NewRoleService(repo persistence.RoleRepositoryInterface) RoleServiceInterface {
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

func (s *RoleService) CreateRole(name, description string) error {
	role := &model.Role{
		Name:        name,
		Description: description,
	}
	return s.repo.CreateRole(role)
}

func (s *RoleService) UpdateRole(id uint, name, description string) error {
	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		return err
	}
	role.Name = name
	role.Description = description
	return s.repo.UpdateRole(role)
}

func (s *RoleService) DeleteRole(id uint) error {
	return s.repo.DeleteRole(id)
}

func (s *RoleService) SeedDefaultPermissions() error {
	perms := []model.Permission{
		{ID: "users_read", Module: "Users", Action: "Read", Description: "View user list and details"},
		{ID: "users_write", Module: "Users", Action: "Write", Description: "Create, update, and delete users"},
		{ID: "roles_read", Module: "Roles", Action: "Read", Description: "View role list and details"},
		{ID: "roles_write", Module: "Roles", Action: "Write", Description: "Create, update, and delete roles"},
		{ID: "news_read", Module: "News", Action: "Read", Description: "View news articles and categories"},
		{ID: "news_write", Module: "News", Action: "Write", Description: "Create, update, and delete news"},
		{ID: "orders_read", Module: "Orders", Action: "Read", Description: "View order list and details"},
		{ID: "orders_write", Module: "Orders", Action: "Write", Description: "Create, update, and delete orders"},
		{ID: "payments_read", Module: "Payments", Action: "Read", Description: "View payment list and details"},
		{ID: "payments_write", Module: "Payments", Action: "Write", Description: "Create, update, and delete payments"},
		{ID: "taxes_read", Module: "Taxes", Action: "Read", Description: "View tax declarations"},
		{ID: "taxes_write", Module: "Taxes", Action: "Write", Description: "Manage tax declarations"},
		{ID: "logs_read", Module: "Logs", Action: "Read", Description: "View system logs"},
	}

	err := s.repo.SeedPermissions(perms)
	if err != nil {
		return err
	}

	// Optionally seed an admin role with all permissions if no roles exist
	roles, err := s.repo.GetAllRoles()
	if err == nil && len(roles) == 0 {
		adminRole := &model.Role{
			Name:        "Admin",
			Description: "System Administrator",
			Permissions: perms,
		}
		_ = s.repo.CreateRole(adminRole)
	}

	return nil
}
