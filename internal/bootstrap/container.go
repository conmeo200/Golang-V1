package bootstrap

import (
	"log"

	"github.com/conmeo200/Golang-V1/internal/infrastructure/rabbitmq"
	"gorm.io/gorm"
)

type Container struct {
	Config   *Config
	DB       *gorm.DB
	RabbitMQ *rabbitmq.RabbitMQ
}

func InitContainer() (*Container, error) {
	cfg 	 := LoadConfig()
	db 		 := InitDatabase(cfg)
	rmq, err := rabbitmq.NewRabbitMQ(cfg.RabbitMQUser, cfg.RabbitMQPassword, cfg.RabbitMQHost, cfg.RabbitMQPort)

	if err != nil {
		return nil, err
	}

	return &Container{
		Config:   cfg,
		DB:       db,
		RabbitMQ: rmq,
	}, nil
}

func (c *Container) Close() {
	if c.RabbitMQ != nil {
		c.RabbitMQ.Close()
	}
	
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
	log.Println("Container resources closed")
}
