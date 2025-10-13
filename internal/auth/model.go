package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

// Claims is a custom struct that embeds jwt.RegisteredClaims and adds custom fields.
type Claims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

type UserInfo struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
