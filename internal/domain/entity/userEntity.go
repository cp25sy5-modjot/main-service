package entity

import "time"

type UserStatus string

const (
	StatusPreActive UserStatus = "pre-active"
	StatusActive    UserStatus = "active"
)

type User struct {
	UserID      string      `gorm:"primaryKey;autoIncrement:false"`
	UserBinding UserBinding `gorm:"embedded"`
	Name        string      `gorm:"length:15"`
	Status      UserStatus `gorm:"default:'pre-active'"`
	Onboarding  bool       `gorm:"default:false"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
}

type UserBinding struct {
	GoogleID   string
	// FacebookID string
	// AppleID    string
}
