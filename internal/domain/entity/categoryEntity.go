package entity

import (
	"time"
)

// category.go
type Category struct {
	CategoryID   string `gorm:"primaryKey;autoIncrement:false"`
	UserID       string `gorm:"index"` 
	CategoryName string `gorm:"length:20"`
	Budget       float64
	ColorCode    string    `gorm:"length:7"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
