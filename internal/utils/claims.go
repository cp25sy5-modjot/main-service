package utils

import (
	"modjot/internal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	r "modjot/internal/response"
)

func GetUserIDFromClaims(c *fiber.Ctx) (string, error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(*auth.Claims)
	if claims == nil || claims.Subject == "" {
		return "", r.InternalServerError(c, "Failed to get user ID from claims")
	}
	return claims.Subject, nil
}
