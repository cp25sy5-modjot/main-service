package auth

import (
	"github.com/cp25sy5-modjot/main-service/internal/shared/config"
	internaljwt "github.com/cp25sy5-modjot/main-service/internal/jwt"
	r "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// RefreshHandler validates a refresh token and issues a new access token.
func RefreshHandler(c *fiber.Ctx, config *config.Auth) error {
	var req RefreshRequest
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	// Parse and validate the refresh token claims.
	token, err := jwt.ParseWithClaims(req.RefreshToken, &internaljwt.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.RefreshTokenSecret), nil
	})

	if err != nil || !token.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	claims, ok := token.Claims.(*internaljwt.Claims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	userInfo := &internaljwt.UserInfo{
		UserID: claims.Subject,
		Name:   claims.Name,
	}
	// Generate a new access token only.
	newAccessToken, _, err := internaljwt.GenerateTokens(userInfo, config)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	return r.OK(c, &TokenResponse{
		AccessToken: newAccessToken,
	}, "Token refreshed successfully")
}
