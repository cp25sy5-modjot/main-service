package database

import (
	"modjot/internal/transaction"

	"gorm.io/gorm"
)

// AutoMigrate runs GORM's automigrate for all entities
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&transaction.Transaction{},
	)
}
