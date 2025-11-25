package model

import (
	"time"

	entity "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
)

type UserInsertReq struct {
	UserBinding entity.UserBinding `json:"user_binding"`
	Name        string             `json:"name" validate:"required,min=2,max=15"`
	DOB         time.Time          `json:"dob"`
}

type UserUpdateReq struct {
	Name string    `json:"name" validate:"min=2,max=15"`
	DOB  time.Time `json:"dob"`
}

type UserRes struct {
	UserBinding entity.UserBinding `json:"user_binding"`
	Name        string             `json:"name"`
	DOB         time.Time          `json:"dob"`
	Status      entity.UserStatus  `json:"status"`
	Onboarding  bool               `json:"onboarding"`
	CreatedAt   time.Time          `json:"created_at"`
}
