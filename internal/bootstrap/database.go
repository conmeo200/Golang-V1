package bootstrap

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	//"gorm.io/plugin/dbresolver"
)

func NewPostgres(cfg *Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitDatabase(cfg *Config) *gorm.DB {
	db, err := NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}
