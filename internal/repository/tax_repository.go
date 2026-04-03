package repository

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaxRepo interface {
	CreateDeclaration(ctx context.Context, decl *model.TaxDeclaration) error
	GetDeclaration(ctx context.Context, id uuid.UUID) (*model.TaxDeclaration, error)
	ListDeclarationsByUserID(ctx context.Context, userID uuid.UUID) ([]model.TaxDeclaration, error)
	
	CreateDependent(ctx context.Context, dep *model.Dependent) error
	ListDependentsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Dependent, error)
}

type TaxRepository struct {
	db *gorm.DB
}

func NewTaxRepository(db *gorm.DB) *TaxRepository {
	return &TaxRepository{db: db}
}

func (r *TaxRepository) CreateDeclaration(ctx context.Context, decl *model.TaxDeclaration) error {
	return r.db.WithContext(ctx).Create(decl).Error
}

func (r *TaxRepository) GetDeclaration(ctx context.Context, id uuid.UUID) (*model.TaxDeclaration, error) {
	var decl model.TaxDeclaration
	err := r.db.WithContext(ctx).Preload("Incomes").First(&decl, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &decl, nil
}

func (r *TaxRepository) ListDeclarationsByUserID(ctx context.Context, userID uuid.UUID) ([]model.TaxDeclaration, error) {
	var decls []model.TaxDeclaration
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&decls).Error
	return decls, err
}

func (r *TaxRepository) CreateDependent(ctx context.Context, dep *model.Dependent) error {
	return r.db.WithContext(ctx).Create(dep).Error
}

func (r *TaxRepository) ListDependentsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Dependent, error) {
	var deps []model.Dependent
	err := r.db.WithContext(ctx).Where("user_id = ? AND (active_to IS NULL OR active_to > NOW())", userID).Find(&deps).Error
	return deps, err
}
