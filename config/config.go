package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgreSQL *PostgreSQL
	App        *Fiber
}

type Fiber struct {
	Host string
	Port string
}

// Database
type PostgreSQL struct {
	Host     string
	Port     string
	Protocol string
	Username string
	Password string
	Database string
	SSLMode  string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	return &Config{
		App: &Fiber{
			Host: getEnv("FIBER_HOST", "localhost"),
			Port: getEnv("FIBER_PORT", "8080"),
		},
		PostgreSQL: &PostgreSQL{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Protocol: getEnv("DB_PROTOCOL", "tcp"),
			Username: getEnv("DB_USERNAME", "appuser"),
			Password: getEnv("DB_PASSWORD", "apppass"),
			Database: getEnv("DB_DATABASE", "modjot"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
