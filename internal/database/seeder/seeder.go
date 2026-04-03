package seeder

import (
	"log"
	"time"

	"github.com/conmeo200/Golang-V1/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	log.Println("Running database seeders...")
	err := seedRolesAndPermissions(db)
	if err != nil {
		log.Printf("Error seeding roles and permissions: %v", err)
		return err
	}

	if err := seedUsers(db); err != nil {
		log.Printf("Error seeding users: %v", err)
		return err
	}

	if err := seedTaxData(db); err != nil {
		log.Printf("Error seeding tax data: %v", err)
		return err
	}

	log.Println("Database seeders completed.")
	return nil
}

func seedUsers(db *gorm.DB) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	now := time.Now().Unix()

	admin := model.User{
		Email:        "admin@example.com",
		PasswordHash: string(hash),
		Role:         "admin", // Legacy string role
		RoleID:       1,       // New Role-based ID (Admin)
		Status:       "active",
		Balance:      1000,
		CreatedAt:    now,
	}

	return db.Where(model.User{Email: admin.Email}).Attrs(admin).FirstOrCreate(&admin).Error
}


func seedTaxData(db *gorm.DB) error {
	var admin model.User
	if err := db.Where("email = ?", "admin@example.com").First(&admin).Error; err != nil {
		log.Printf("Admin user not found for seeding tax data, skipping...")
		return nil
	}

	// 1. Seed Dependents
	dependents := []model.Dependent{
		{UserID: admin.ID, Name: "Nguyễn Văn A", IDNumber: "001201001234", Relationship: "Con", ActiveFrom: time.Now()},
		{UserID: admin.ID, Name: "Trần Thị B", IDNumber: "001201005678", Relationship: "Vợ", ActiveFrom: time.Now()},
	}

	for _, d := range dependents {
		db.FirstOrCreate(&model.Dependent{}, model.Dependent{IDNumber: d.IDNumber}).Updates(d)
	}

	// 2. Seed Declarations
	decls := []model.TaxDeclaration{
		{
			UserID:      admin.ID,
			Type:        model.TaxTypePIT,
			Period:      "2026-M01",
			TotalIncome: 50000000,
			Status:      model.StatusSubmitted,
			TaxPayable:  6550000, // Roughly calculated for 50M with 2 deps (11 + 4.4*2 = 19.8M deduction, taxable 30.2M)
		},
		{
			UserID:      admin.ID,
			Type:        model.TaxTypeHousehold,
			Period:      "2026-Q1",
			TotalIncome: 120000000,
			Status:      model.StatusApproved,
			TaxPayable:  1800000, // 1.5% of 120M
		},
	}

	for _, d := range decls {
		db.FirstOrCreate(&model.TaxDeclaration{}, model.TaxDeclaration{UserID: d.UserID, Period: d.Period}).Updates(d)
	}

	return nil
}

func seedRolesAndPermissions(db *gorm.DB) error {
	permissions := []model.Permission{
		// User Management
		{ID: "usr_list", Module: "User Management", Action: "List", Description: "Can view user list"},
		{ID: "usr_create", Module: "User Management", Action: "Create", Description: "Can create new users"},
		{ID: "usr_update", Module: "User Management", Action: "Update", Description: "Can edit users"},
		{ID: "usr_delete", Module: "User Management", Action: "Delete", Description: "Can delete users"},

		// Role Management
		{ID: "role_list", Module: "Role Management", Action: "List", Description: "Can view role list"},
		{ID: "role_create", Module: "Role Management", Action: "Create", Description: "Can create new roles"},
		{ID: "role_update", Module: "Role Management", Action: "Update", Description: "Can edit roles"},
		{ID: "role_delete", Module: "Role Management", Action: "Delete", Description: "Can delete roles"},

		// News Management
		{ID: "news_list", Module: "News Management", Action: "List", Description: "Can view news"},
		{ID: "news_create", Module: "News Management", Action: "Create", Description: "Can create news"},
		{ID: "news_update", Module: "News Management", Action: "Update", Description: "Can edit news"},
		{ID: "news_delete", Module: "News Management", Action: "Delete", Description: "Can delete news"},
	}

	for _, p := range permissions {
		err := db.FirstOrCreate(&p, model.Permission{ID: p.ID}).Error
		if err != nil {
			return err
		}
	}

	roles := []model.Role{
		{
			ID:          1,
			Name:        "Admin",
			Description: "Full access to all modules and configurations",
		},
		{
			ID:          2,
			Name:        "Editor",
			Description: "Can edit content, news, and users",
		},
		{
			ID:          3,
			Name:        "Viewer",
			Description: "Read-only access to specific dashboards",
		},
	}

	for _, r := range roles {
		err := db.FirstOrCreate(&model.Role{}, model.Role{ID: r.ID}).Updates(r).Error
		if err != nil {
			return err
		}
	}

	// Assign permissions
	var adminRole, editorRole, viewerRole model.Role
	db.First(&adminRole, 1)
	db.First(&editorRole, 2)
	db.First(&viewerRole, 3)

	var allPerms []model.Permission
	db.Find(&allPerms)
	db.Model(&adminRole).Association("Permissions").Replace(allPerms)

	var editorPerms []model.Permission
	db.Where("id IN ?", []string{"usr_list", "usr_update", "news_list", "news_create", "news_update"}).Find(&editorPerms)
	db.Model(&editorRole).Association("Permissions").Replace(editorPerms)

	var viewerPerms []model.Permission
	db.Where("id IN ?", []string{"usr_list", "role_list", "news_list"}).Find(&viewerPerms)
	db.Model(&viewerRole).Association("Permissions").Replace(viewerPerms)

	return nil
}
