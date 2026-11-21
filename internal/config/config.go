package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DBHost        string
	DBPort        string
	DBName        string
	DBUser        string
	DBPass        string
	JWTSecret     string
	AllowedOrigin string
}

func Load() (*Config, error) {
	_ = godotenv.Load() // ignore error if .env doesn't exist

	cfg := &Config{
		Port:          getEnv("PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBName:        getEnv("DB_NAME", "restobill"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPass:        getEnv("DB_PASSWORD", "postgres"),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "http://localhost:5173"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
