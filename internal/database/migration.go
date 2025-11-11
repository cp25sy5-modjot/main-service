package database

import (
	userModel "github.com/cp25sy5-modjot/main-service/internal/user/model"
	categoryModel "github.com/cp25sy5-modjot/main-service/internal/category/model"
	transactionModel "github.com/cp25sy5-modjot/main-service/internal/transaction/model"

	"gorm.io/gorm"
)

// AutoMigrate runs GORM's automigrate for all entities
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&transactionModel.Transaction{},
		&userModel.User{},
		&categoryModel.Category{},
	)
}
