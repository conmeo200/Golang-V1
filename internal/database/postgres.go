package database

import (
	"fmt"

	"github.com/conmeo200/Golang-V1/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	//"gorm.io/plugin/dbresolver"
)

func NewPostgres(cfg *config.Config) (*gorm.DB, error) {

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
