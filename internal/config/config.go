package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port             string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecretKey     string
	DBReplicaHosts   []string
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUser        string
	RabbitMQPassword    string

	// Stripe Configuration
	StripePublicKey     string
	StripeSecretKey     string
	StripeWebhookSecret string

	// PayPal Configuration
	PayPalClientID      string
	PayPalSecret        string
	PayPalEnvironment   string // "sandbox" or "live"
}

func Load() *Config {
	v := viper.New()

	// Thiết lập các giá trị mặc định
	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")

	// Đọc từ file .env nếu tồn tại
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	if err := v.ReadInConfig(); err != nil {
		log.Printf("Note: .env file not found or could not be read: %s", err)
	}

	// Tự động đọc từ biến môi trường hệ thống
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return &Config{
		Port:         v.GetString("APP_PORT"),
		DBHost:       v.GetString("DB_HOST"),
		DBPort:       v.GetString("DB_PORT"),
		DBUser:       v.GetString("DB_USER"),
		DBPassword:   v.GetString("DB_PASSWORD"),
		DBName:       v.GetString("DB_NAME"),
		JWTSecretKey: v.GetString("JWT_SECRET_KEY"),
		DBReplicaHosts: func() []string {
			hosts := v.GetString("DB_REPLICA_HOSTS")
			if hosts == "" {
				return []string{}
			}
			return strings.Split(hosts, ",")
		}(),
		RabbitMQHost:        v.GetString("RABBITMQ_HOST"),
		RabbitMQPort:        v.GetString("RABBITMQ_PORT"),
		RabbitMQUser:        v.GetString("RABBITMQ_USER"),
		RabbitMQPassword:    v.GetString("RABBITMQ_PASSWORD"),
		
		StripePublicKey:     v.GetString("STRIPE_PUBLIC_KEY"),
		StripeSecretKey:     v.GetString("STRIPE_SECRET_KEY"),
		StripeWebhookSecret: v.GetString("STRIPE_WEBHOOK_SECRET"),
		
		PayPalClientID:      v.GetString("PAYPAL_CLIENT_ID"),
		PayPalSecret:        v.GetString("PAYPAL_SECRET"),
		PayPalEnvironment:   v.GetString("PAYPAL_ENVIRONMENT"),
	}
}
