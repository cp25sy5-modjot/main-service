package database

import (
	"github.com/cp25sy5-modjot/main-service/internal/transaction"
	"github.com/cp25sy5-modjot/main-service/internal/user"

	"gorm.io/gorm"
)

// AutoMigrate runs GORM's automigrate for all entities
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&transaction.Transaction{},
		&user.User{},
	)
}
