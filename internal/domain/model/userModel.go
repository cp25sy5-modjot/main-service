package model

import (
	"time"

	entity "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
)

type UserStatus string

const (
	UserStatusPreActive UserStatus = "pre-active"
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
)

type UserInsertReq struct {
	UserBinding entity.UserBinding `json:"user_binding"`
	Name        string             `json:"name" validate:"required,min=1,max=15"`
}

type UserUpdateReq struct {
	Name string `json:"name" validate:"min=1,max=15"`
}

type UserRes struct {
	UserBinding UserBinding `json:"user_binding"`
	Name        string      `json:"name"`
	Status      string      `json:"status"`
	Onboarding  bool        `json:"onboarding"`
	CreatedAt   time.Time   `json:"created_at"`
}

type UserBinding struct {
	GoogleID string `json:"google_id"`
	// FacebookID string `json:"facebook_id"`
	// AppleID    string `json:"apple_id"`
}
