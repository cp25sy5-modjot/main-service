package auth

import (
	"github.com/cp25sy5-modjot/main-service/internal/shared/config"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	r "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"

	u "github.com/cp25sy5-modjot/main-service/internal/user/service"
	c "github.com/cp25sy5-modjot/main-service/internal/category/service"
	"github.com/gofiber/fiber/v2"
)

// MockLoginHandler handles mock login requests for testing purposes Only in non-production environments.
func MockLoginHandler(c *fiber.Ctx, usvc u.Service, csvc c.Service, config *config.Auth) error {
	userName := c.FormValue("userName")

	if userName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "userName is required")
	}

	user, err := usvc.GetByID(userName)
	if err != nil {
		user, err = usvc.CreateMockUser(&u.UserCreateInput{
			Name: userName,
		}, userName)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
		}
		// Create default categories for the new mock user
		if err := csvc.CreateDefaultCategories(user.UserID); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to create default categories for mock user")
		}
	}

	userInfo := &jwt.UserInfo{
		UserID: user.UserID,
		Name:   user.Name,
	}

	accessToken, refreshToken, err := jwt.GenerateTokens(userInfo, config)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate tokens")
	}
	return r.OK(c, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, "Login successful")
}
