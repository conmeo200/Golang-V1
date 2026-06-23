package bootstrap

import (
	"log"

	rabbitmqinfra "github.com/conmeo200/Golang-V1/internal/infrastructure/rabbitmq"
)

func RunAPI() error {
	cfg := LoadConfig()
	db := InitDatabase(cfg)

	rabbitConn, err := rabbitmqinfra.NewRabbitMQ(cfg.RabbitMQUser, cfg.RabbitMQPassword, cfg.RabbitMQHost, cfg.RabbitMQPort)
	if err != nil {
		return err
	}

	container := BuildContainer(
		cfg,
		db,
		rabbitConn,
	)

	app := NewAPIApp(container)
	
	log.Printf("🚀 Web Server starting on %s", app.Server.Addr)

	return app.Server.ListenAndServe()
}
