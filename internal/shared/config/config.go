package config

import (
	"os"

	"github.com/joho/godotenv"
)

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
		Auth: &Auth{
			AccessTokenSecret:  os.Getenv("ACCESS_TOKEN_SECRET"),
			RefreshTokenSecret: os.Getenv("REFRESH_TOKEN_SECRET"),
			AccessTokenTTL:     os.Getenv("ACCESS_TOKEN_TTL"),
			RefreshTokenTTL:    os.Getenv("REFRESH_TOKEN_TTL"),
			Issuer:             os.Getenv("JWT_ISSUER"),
		},
		Google: &Google{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		},
		AIService: &AIService{
			Url: os.Getenv("AI_WRAPPER_ADDR"),
		},
		Redis: &Redis{
			Addr: os.Getenv("REDIS_ADDR"),
		},
		UploadDir: os.Getenv("UPLOAD_DIR"),
	}
}

