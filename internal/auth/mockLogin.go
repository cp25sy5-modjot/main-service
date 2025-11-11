package auth

import (
	"github.com/cp25sy5-modjot/main-service/internal/config"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	r "github.com/cp25sy5-modjot/main-service/internal/response/success"

	u "github.com/cp25sy5-modjot/main-service/internal/user"
	"github.com/gofiber/fiber/v2"
)

func MockLoginHandler(c *fiber.Ctx, service *u.Service, config *config.Auth) error {
	userName := c.FormValue("userName")

	if userName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "userName is required")
	}

	user, err := service.Create(&u.UserInsertReq{
		Email: userName + "@mock.com",
		Name:  userName,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
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
