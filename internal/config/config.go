package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port         string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	JWTSecretKey string
}

func Load() *Config {
	v := viper.New()

	// Thiết lập các giá trị mặc định
	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")

	// Tự động đọc từ biến môi trường
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Đọc từ file nếu tồn tại
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Error reading config file: %s", err)
		}
	}

	return &Config{
		Port:         v.GetString("APP_PORT"),
		DBHost:       v.GetString("DB_HOST"),
		DBPort:       v.GetString("DB_PORT"),
		DBUser:       v.GetString("DB_USER"),
		DBPassword:   v.GetString("DB_PASSWORD"),
		DBName:       v.GetString("DB_NAME"),
		JWTSecretKey: v.GetString("JWT_SECRET_KEY"),
	}
}