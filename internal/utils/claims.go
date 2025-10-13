package utils

import (
	"modjot/internal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserIDFromClaims(c *fiber.Ctx) string {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(*auth.Claims)
	if claims == nil || claims.Subject == "" {
		return ""
	}
	return claims.Subject
}
