package auth

import (
	"github.com/cp25sy5-modjot/main-service/internal/config"
	r "github.com/cp25sy5-modjot/main-service/internal/response/success"

	"github.com/gofiber/fiber/v2"
)

func MockLoginHandler(c *fiber.Ctx, config *config.Auth) error {
	userID := c.FormValue("userID")
	userName := c.FormValue("userName")

	if userID == "" || userName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "userID and userName are required")
	}

	// For a real app, userID would come from your database.

	userInfo := &UserInfo{
		UserID: userID,
		Name:   userName,
	}
	// Generate both access and refresh tokens.
	accessToken, refreshToken, err := GenerateTokens(userInfo, config)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate tokens")
	}
	return r.OK(c, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, "Login successful")
}
