package auth

import (
	"modjot/internal/config"
	r "modjot/internal/response"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)


// RefreshHandler validates a refresh token and issues a new access token.
func RefreshHandler(c *fiber.Ctx, config *config.Auth) error {
	req := new(RefreshRequest)
	if err := c.BodyParser(req); err != nil {
		return r.BadRequest(c, "Invalid JSON body")
	}

	// Parse and validate the refresh token.
	token, err := jwt.ParseWithClaims(req.RefreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.RefreshTokenSecret), nil
	})

	if err != nil || !token.Valid {
		return r.Unauthorized(c, "Invalid or expired refresh token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return r.InternalServerError(c, "Failed to parse claims")
	}

	userInfo := &UserInfo{
		UserID: claims.Subject,
		Name:   claims.Name,
	}
	// Generate a new access token only.
	newAccessToken, _, err := GenerateTokens(userInfo, config)
	if err != nil {
		return r.InternalServerError(c, "Failed to generate new access token")
	}

	return r.OK(c, fiber.Map{
		"access_token": newAccessToken,
	})
}
