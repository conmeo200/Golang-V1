package model

type Permission struct {
	ID          string `gorm:"size:50;primaryKey"` // e.g. "usr_list"
	Module      string `gorm:"size:100;not null"`
	Action      string `gorm:"size:100;not null"`
	Description string `gorm:"size:255"`

	// Relationships
	Roles []Role `gorm:"many2many:role_permissions;"`
}
