package repository

import (
	"github.com/conmeo200/Golang-V1/internal/model"
	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	GetAllRoles() ([]model.Role, error)
	GetRoleByID(id uint) (*model.Role, error)
	GetRoleWithPermissions(id uint) (*model.Role, error)
	GetAllPermissions() ([]model.Permission, error)
	UpdateRolePermissions(roleID uint, permissionIDs []string) error
	GetUsersCountByRoleName(roleName string) (int64, error)
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepositoryInterface {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetAllRoles() ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) GetRoleByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetRoleWithPermissions(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetAllPermissions() ([]model.Permission, error) {
	var perms []model.Permission
	err := r.db.Find(&perms).Error
	return perms, err
}

func (r *RoleRepository) UpdateRolePermissions(roleID uint, permissionIDs []string) error {
	var role model.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	var perms []model.Permission
	if len(permissionIDs) > 0 {
		if err := r.db.Where("id IN ?", permissionIDs).Find(&perms).Error; err != nil {
			return err
		}
	}

	return r.db.Model(&role).Association("Permissions").Replace(perms)
}

func (r *RoleRepository) GetUsersCountByRoleName(roleName string) (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("role = ?", roleName).Count(&count).Error
	return count, err
}
