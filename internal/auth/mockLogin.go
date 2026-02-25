package auth

import (
	"errors"
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	"github.com/cp25sy5-modjot/main-service/internal/shared/config"
	r "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
	"gorm.io/gorm"

	c "github.com/cp25sy5-modjot/main-service/internal/category/service"
	utils "github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	u "github.com/cp25sy5-modjot/main-service/internal/user/service"
	"github.com/gofiber/fiber/v2"
)

// MockLoginHandler handles mock login requests for testing purposes Only in non-production environments.
func MockLoginHandler(c *fiber.Ctx, usvc u.Service, csvc c.Service, config *config.Auth) error {
	userName := c.FormValue("userName")

	if userName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "userName is required")
	}

	user, err := usvc.GetByID(userName)

	if user != nil {
		if user.Status == e.UserStatusInactive {
			expireAt := user.UnsubscribedAt.Add(30 * 24 * time.Hour)
			remaining := time.Until(expireAt)

			return utils.AccountDeactivated(int(remaining.Seconds()))
		}

	}

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user")
		}
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

func MockRestoreHandler(c *fiber.Ctx, usvc u.Service, config *config.Auth) error {
	userName := c.FormValue("userName")

	if userName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "userName is required")
	}

	user, err := usvc.RestoreByUserID(userName)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to restore user")
	}

	userInfo := &jwt.UserInfo{
		UserID: user.UserID,
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
