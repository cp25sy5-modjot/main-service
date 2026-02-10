package entity

import "time"

type FavoriteItem struct {
	FavoriteID string `gorm:"primaryKey;autoIncrement:false"`
	UserID     string
	Title      string
	Price      float64
	CategoryID string
	Position   int
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	// Relationships
	Category Category `gorm:"foreignKey:CategoryID;references:CategoryID"`
}
