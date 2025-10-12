package config

import (
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
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load()
	}
	return &Config{
		App: &Fiber{
			Host: os.Getenv("FIBER_HOST"),
			Port: os.Getenv("FIBER_PORT"),
		},
		PostgreSQL: &PostgreSQL{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Protocol: os.Getenv("POSTGRES_PROTOCOL"),
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DB"),
			SSLMode:  os.Getenv("POSTGRES_SSL_MODE"),
		},
	}
}
