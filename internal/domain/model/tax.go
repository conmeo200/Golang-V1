package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TaxType represents the category of the tax filing
type TaxType string

const (
	TaxTypePIT       TaxType = "PIT"       // Personal Income Tax (Individual)
	TaxTypeHousehold TaxType = "HOUSEHOLD" // Household Business Tax
	TaxTypeCIT       TaxType = "CIT"       // Corporate Income Tax
)

// TaxStatus represents the state of the declaration
type TaxStatus string

const (
	StatusDraft     TaxStatus = "DRAFT"
	StatusSubmitted TaxStatus = "SUBMITTED"
	StatusApproved  TaxStatus = "APPROVED"
	StatusRejected  TaxStatus = "REJECTED"
)

// TaxDeclaration is the main entity for a tax filing
type TaxDeclaration struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`
	Type           TaxType        `gorm:"type:varchar(20);not null" json:"type"`
	Period         string         `gorm:"type:varchar(20);not null" json:"period"` // e.g., "2026-Q1", "2026-M03"
	TotalIncome    int64          `gorm:"not null" json:"total_income"`
	TotalDeduction int64          `json:"total_deduction"`
	TaxPayable     int64          `json:"tax_payable"`
	Status         TaxStatus      `gorm:"type:varchar(20);default:'DRAFT'" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Incomes []IncomeRecord `gorm:"foreignKey:DeclarationID" json:"incomes"`
}

func (t *TaxDeclaration) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// IncomeRecord stores individual income items
type IncomeRecord struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	DeclarationID uuid.UUID      `gorm:"type:uuid;index" json:"declaration_id"`
	SourceName    string         `gorm:"type:varchar(255);not null" json:"source_name"`
	Amount        int64          `gorm:"not null" json:"amount"`
	ReceivedAt    time.Time      `json:"received_at"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// Dependent manages deductions for family circumstances
type Dependent struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	IDNumber     string         `gorm:"type:varchar(50)" json:"id_number"`
	Relationship string         `gorm:"type:varchar(50)" json:"relationship"`
	ActiveFrom   time.Time      `json:"active_from"`
	ActiveTo     *time.Time     `json:"active_to"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

