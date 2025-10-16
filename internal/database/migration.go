package database

import (
	"github.com/cp25sy5-modjot/main-service/internal/transaction"

	"gorm.io/gorm"
)

// AutoMigrate runs GORM's automigrate for all entities
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&transaction.Transaction{},
	)
}
