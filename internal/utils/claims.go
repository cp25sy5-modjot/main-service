package utils

import (
	"github.com/cp25sy5-modjot/main-service/internal/auth"

	r "github.com/cp25sy5-modjot/main-service/internal/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserIDFromClaims(c *fiber.Ctx) (string, error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(*auth.Claims)
	if claims == nil || claims.Subject == "" {
		return "", r.InternalServerError(c, "Failed to get user ID from claims")
	}
	return claims.Subject, nil
}
