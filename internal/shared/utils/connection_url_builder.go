package utils

import (
	"fmt"

	"github.com/cp25sy5-modjot/main-service/internal/shared/config"
)

func PostgresUrlBuilder(cfg *config.Config) (string, error) {
	url := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.Username,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Database,
		cfg.PostgreSQL.SSLMode,
	)
	return url, nil
}

func AppUrlBuilder(cfg *config.Config) (string, error) {
	url := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	return url, nil
}
