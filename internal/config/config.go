package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgreSQL *PostgreSQL
	App        *Fiber
	Auth       *Auth
	Google     *Google
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
		Auth: &Auth{
			AccessTokenExpireMinutes:  parseEnvToInt64("JWT_ACCESS_TOKEN_EXPIRE_MINUTES"),
			RefreshTokenExpireMinutes: parseEnvToInt64("JWT_REFRESH_TOKEN_EXPIRE_MINUTES"),
			SecretKey:                 os.Getenv("JWT_SECRET"),
			Issuer:                    os.Getenv("JWT_ISSUER"),
		},
		Google: &Google{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		},
	}
}

func parseEnvToInt64(key string) int64 {
	valueStr := os.Getenv(key)
	var value int64
	fmt.Sscan(valueStr, &value)
	return value
}
