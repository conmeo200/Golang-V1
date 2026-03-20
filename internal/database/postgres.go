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

	// Cấu hình Read Replicas nếu có
	// if len(cfg.DBReplicaHosts) > 0 {
	// 	var replicas []gorm.Dialector
	// 	for _, host := range cfg.DBReplicaHosts {
	// 		replicaDsn := fmt.Sprintf(
	// 			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	// 			host,
	// 			cfg.DBUser,
	// 			cfg.DBPassword,
	// 			cfg.DBName,
	// 			cfg.DBPort,
	// 		)
	// 		replicas = append(replicas, postgres.Open(replicaDsn))
	// 	}

	// 	err = db.Use(dbresolver.Register(dbresolver.Config{
	// 		Sources:  []gorm.Dialector{postgres.Open(dsn)},
	// 		Replicas: replicas,
	// 		Policy:   dbresolver.RandomPolicy{},
	// 	}))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return db, nil
}
