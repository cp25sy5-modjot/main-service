package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusPreActive UserStatus = "pre-active"
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
)

type User struct {
	UserID         string         `gorm:"primaryKey;autoIncrement:false"`
	UserBinding    UserBinding    `gorm:"embedded"`
	Name           string         `gorm:"length:15"`
	Status         UserStatus     `gorm:"default:'pre-active'"`
	Onboarding     bool           `gorm:"default:false"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	UnsubscribedAt *time.Time     `gorm:"default:null"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type UserBinding struct {
	GoogleID string
	// FacebookID string
	// AppleID    string
}
