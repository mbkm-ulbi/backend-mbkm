package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName  string
	AppPort  string
	AppEnv   string
	DBHost   string
	DBPort   string
	DBName   string
	DBUser   string
	DBPass   string
	JWTSecret string
	JWTExpiry time.Duration
}

var AppConfig *Config

func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		// Not a fatal error, env vars might be set directly
	}

	expiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))

	AppConfig = &Config{
		AppName:   getEnv("APP_NAME", "mbkm-go"),
		AppPort:   getEnv("APP_PORT", "3000"),
		AppEnv:    getEnv("APP_ENV", "development"),
		DBHost:    getEnv("DB_HOST", "127.0.0.1"),
		DBPort:    getEnv("DB_PORT", "3306"),
		DBName:    getEnv("DB_DATABASE", "mbkm"),
		DBUser:    getEnv("DB_USERNAME", "root"),
		DBPass:    getEnv("DB_PASSWORD", ""),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiry: expiry,
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
