package user

import "time"

type UserStatus string

const (
	StatusPreActive UserStatus = "pre-active"
	StatusActive    UserStatus = "active"
)

type User struct {
	UserID      string      `gorm:"primaryKey;autoIncrement:false" json:"user_id"`
	UserBinding UserBinding `gorm:"embedded" json:"user_binding"`
	Name        string      `gorm:"length:15" json:"name"`
	DOB         time.Time   `json:"dob"`
	Status      UserStatus  `gorm:"default:'pre-active'" json:"status"`
	Onboarding  bool        `gorm:"default:false" json:"onboarding"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

type UserBinding struct {
	GoogleID   string `json:"google_id"`
	FacebookID string `json:"facebook_id"`
	AppleID    string `json:"apple_id"`
}
