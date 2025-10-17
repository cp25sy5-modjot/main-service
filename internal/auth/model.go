package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type GoogleTokenRequest struct {
	Code         string `json:"code" validate:"required"`
	CodeVerifier string `json:"code_verifier" validate:"required,min=43,max=128"`
}

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
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
