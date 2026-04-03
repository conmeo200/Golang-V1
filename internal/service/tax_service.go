package service

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/repository"
	"github.com/google/uuid"
)

// TaxServiceInterface defines the core calculation and management for taxes
type TaxServiceInterface interface {
	CalculatePIT(taxableIncome int64) int64
	CalculateHouseholdTax(revenue int64) int64
	GetBaseDeduction() int64
	GetDependentDeduction() int64
	ListDeclarations(ctx context.Context, userID uuid.UUID) ([]model.TaxDeclaration, error)
	CreateDeclaration(ctx context.Context, decl *model.TaxDeclaration) error
}

type TaxService struct {
	repo repository.TaxRepo
}

func NewTaxService(repo repository.TaxRepo) *TaxService {
	return &TaxService{repo: repo}
}

func (s *TaxService) ListDeclarations(ctx context.Context, userID uuid.UUID) ([]model.TaxDeclaration, error) {
	return s.repo.ListDeclarationsByUserID(ctx, userID)
}

func (s *TaxService) CreateDeclaration(ctx context.Context, decl *model.TaxDeclaration) error {
	// 1. Calculate tax based on type
	if decl.Type == model.TaxTypePIT {
		// Get deductions
		dependents, _ := s.repo.ListDependentsByUserID(ctx, decl.UserID)
		totalDeduction := s.GetBaseDeduction() + int64(len(dependents))*s.GetDependentDeduction()
		decl.TotalDeduction = totalDeduction

		taxable := decl.TotalIncome - totalDeduction
		decl.TaxPayable = s.CalculatePIT(taxable)
	} else if decl.Type == model.TaxTypeHousehold {
		decl.TaxPayable = s.CalculateHouseholdTax(decl.TotalIncome)
	}

	// 2. Save via Repo
	return s.repo.CreateDeclaration(ctx, decl)
}


// CalculatePIT calculates Personal Income Tax based on 2026 progressive rates
// income is the taxable income per month (after all deductions)
func (s *TaxService) CalculatePIT(taxableIncome int64) int64 {
	if taxableIncome <= 0 {
		return 0
	}

	var tax int64

	// Biểu thuế lũy tiến từng phần (Vietnamese PIT 2026 Progressive Rates)
	// Bậc 1: Đến 5 triệu đồng (5%)
	// Bậc 2: Trên 5 đến 10 triệu đồng (10%)
	// Bậc 3: Trên 10 đến 18 triệu đồng (15%)
	// Bậc 4: Trên 18 đến 32 triệu đồng (20%)
	// Bậc 5: Trên 32 đến 52 triệu đồng (25%)
	// Bậc 6: Trên 52 đến 80 triệu đồng (30%)
	// Bậc 7: Trên 80 triệu đồng (35%)

	if taxableIncome <= 5000000 {
		tax = taxableIncome * 5 / 100
	} else if taxableIncome <= 10000000 {
		tax = (taxableIncome-5000000)*10/100 + 250000
	} else if taxableIncome <= 18000000 {
		tax = (taxableIncome-10000000)*15/100 + 750000
	} else if taxableIncome <= 32000000 {
		tax = (taxableIncome-18000000)*20/100 + 1950000
	} else if taxableIncome <= 52000000 {
		tax = (taxableIncome-32000000)*25/100 + 4750000
	} else if taxableIncome <= 80000000 {
		tax = (taxableIncome-52000000)*30/100 + 9750000
	} else {
		tax = (taxableIncome-80000000)*35/100 + 18150000
	}

	return tax
}

// GetBaseDeduction returns the fixed personal deduction (e.g., 11,000,000 VND)
func (s *TaxService) GetBaseDeduction() int64 {
	return 11000000
}

// GetDependentDeduction returns the deduction per dependent (e.g., 4,400,000 VND)
func (s *TaxService) GetDependentDeduction() int64 {
	return 4400000
}

// CalculateHouseholdTax calculates tax for business households
// As a placeholder, we use 1.5% (1.0% VAT + 0.5% PIT) for general trade/sales
func (s *TaxService) CalculateHouseholdTax(revenue int64) int64 {
	if revenue <= 100000000/12 { // Low income threshold
		return 0
	}
	return revenue * 15 / 1000 // 1.5%
}
