package category

import (
	"time"
)

type Category struct {
	CategoryID   string    `gorm:"primaryKey;autoIncrement:false" json:"category_id" validate:"required"`
	UserID       string    `gorm:"primaryKey;autoIncrement:false" json:"user_id" validate:"required"`
	CategoryName string    `gorm:"length:100" json:"category_name" validate:"required"`
	Budget       float64   `json:"budget" validate:"required"`
	ColorCode    string    `gorm:"length:7" json:"color_code" validate:"required"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
