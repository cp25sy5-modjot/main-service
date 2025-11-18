package auth

import (
	"context"
	"log"

	"github.com/cp25sy5-modjot/main-service/internal/config"

	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	r "github.com/cp25sy5-modjot/main-service/internal/response/success"
	u "github.com/cp25sy5-modjot/main-service/internal/user/service"
	userModel "github.com/cp25sy5-modjot/main-service/internal/user/model"
	"github.com/cp25sy5-modjot/main-service/internal/utils"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/idtoken"
)

func HandleGoogleTokenExchange(c *fiber.Ctx, service *u.Service, config *config.Config) error {
	var req GoogleTokenRequest
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	payload, err := validateIDToken(req.IdToken, config.Google)
	if err != nil {
		log.Printf("Error validating ID token: %v", err)
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid ID token")
	}

	userInfo := getUserInfoFromPayload(payload, service)

	accessToken, refreshToken, err := jwt.GenerateTokens(userInfo, config.Auth)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate tokens")
	}

	return r.OK(c, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, "Login successful")
}

func validateIDToken(idToken string, config *config.Google) (*idtoken.Payload, error) {
	ctx := context.Background()
	payload, err := idtoken.Validate(ctx, idToken, config.ClientID)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func getUserInfoFromPayload(payload *idtoken.Payload, service *u.Service) *jwt.UserInfo {
	googleID := payload.Subject

	user, err := service.GetByGoogleID(googleID)
	if err != nil {
		name := payload.Claims["given_name"].(string)
		if name == "" {
			name = payload.Claims["name"].(string)
		}
		if name == "" {
			name = "New User"
		}
		user, err = service.Create(&userModel.UserInsertReq{
			Name:  name,
			UserBinding: userModel.UserBinding{
				GoogleID: googleID,
			},
		})
		if err != nil {
			fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
		}
		log.Printf("Created new user!")
	}

	return &jwt.UserInfo{
		UserID: user.UserID,
		Name:   user.Name,
	}
}
