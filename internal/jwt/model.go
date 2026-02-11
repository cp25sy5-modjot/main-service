package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

// Claims is a custom struct that embeds jwt.RegisteredClaims and adds custom fields.
type Claims struct {
	Type string `json:"type"` // access | refresh
	jwt.RegisteredClaims
}

type UserInfo struct {
	UserID string `json:"user_id"`
}

