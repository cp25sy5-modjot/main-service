package user

import (
	"time"
)

type UserInsertReq struct {
	UserID string    `json:"user_id"`
	Name   string    `json:"name" validate:"required,min=2,max=100"`
	Email  string    `json:"email" validate:"required,email"`
	DOB    time.Time `json:"dob"`
}

type UserUpdateReq struct {
	Name string    `json:"name" validate:"min=2,max=100"`
	DOB  time.Time `json:"dob"`
}

type UserRes struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	DOB       time.Time `json:"dob"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
