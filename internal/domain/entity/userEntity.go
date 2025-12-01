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
	DOB         time.Time
	Status      UserStatus `gorm:"default:'pre-active'"`
	Onboarding  bool       `gorm:"default:false"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`

	// has-many Categories – ลบ user แล้วลบทุก category
	Categories []Category `gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// has-many Transactions – ลบ user แล้วลบทุก transaction
	Transactions []Transaction `gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserBinding struct {
	GoogleID   string
	FacebookID string
	AppleID    string
}
