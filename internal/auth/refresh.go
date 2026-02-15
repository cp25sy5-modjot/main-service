package auth

import (
	"errors"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	internaljwt "github.com/cp25sy5-modjot/main-service/internal/jwt"
	"github.com/cp25sy5-modjot/main-service/internal/shared/config"
	r "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	u "github.com/cp25sy5-modjot/main-service/internal/user/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RefreshHandler(c *fiber.Ctx, usvc u.Service, conf *config.Auth) error {
	var req RefreshRequest
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	token, err := jwt.ParseWithClaims(
		req.RefreshToken,
		&internaljwt.Claims{},
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, fiber.ErrUnauthorized
			}

			return []byte(conf.RefreshTokenSecret), nil
		},
	)

	if errors.Is(err, jwt.ErrTokenExpired) {
		return fiber.NewError(fiber.StatusUnauthorized, "Refresh token expired")
	}

	if err != nil || !token.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	claims, ok := token.Claims.(*internaljwt.Claims)
	if !ok || claims.Type != "refresh" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	user, err := usvc.GetByID(claims.Subject)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}

	if user.Status != e.UserStatusActive {
		return fiber.NewError(fiber.StatusUnauthorized, "User disabled")
	}

	userInfo := &internaljwt.UserInfo{
		UserID: user.UserID,
	}

	newAccessToken, _, err := internaljwt.GenerateTokens(userInfo, conf)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	return r.OK(c, &TokenResponse{
		AccessToken: newAccessToken,
	}, "Token refreshed successfully")
}
