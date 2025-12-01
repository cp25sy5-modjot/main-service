package database

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"

	"gorm.io/gorm"
)

// AutoMigrate runs GORM's automigrate for all entities
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&e.User{},
		&e.Category{},
		&e.Transaction{},
		// &e.TransactionItem{},
	)
}
