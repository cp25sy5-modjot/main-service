package entity

import "time"

type PushToken struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	Token     string `gorm:"uniqueIndex"`
	Platform  string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Relationships
	User User `gorm:"foreignKey:UserID;references:UserID"`
}
