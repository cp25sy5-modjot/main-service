package auth

import (
	"context"
	"errors"
	"log"

	"github.com/cp25sy5-modjot/main-service/internal/shared/config"
	"gorm.io/gorm"

	c "github.com/cp25sy5-modjot/main-service/internal/category/service"
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	r "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	u "github.com/cp25sy5-modjot/main-service/internal/user/service"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/idtoken"
)

func HandleGoogleTokenExchange(c *fiber.Ctx, usvc u.Service, csvc c.Service, config *config.Config) error {
	payload, err := parseGoogleToken(c, config)
	if err != nil {
		return err
	}

	googleID := payload.Subject
	if googleID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "google id not found in token")
	}

	user, err := usvc.GetByGoogleID(googleID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to get user")
		}
	}

	if user != nil {
		if user.Status == e.StatusInactive {
			return fiber.NewError(
				fiber.StatusForbidden,
				"account has been deactivated",
			)
		}
	}

	if user == nil {
		var name string
		if v, ok := payload.Claims["given_name"].(string); ok && v != "" {
			name = v
		} else if v, ok := payload.Claims["name"].(string); ok && v != "" {
			name = v
		} else {
			name = "New User"
		}

		user, err = usvc.Create(&u.UserCreateInput{
			Name: name,
			UserBinding: e.UserBinding{
				GoogleID: googleID,
			},
		})
		if err != nil || user == nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to create user")
		}
		log.Printf("created new user with googleID %s", googleID)

		// Create default categories for the new user
		if err := csvc.CreateDefaultCategories(user.UserID); err != nil {
			log.Printf("failed to create default categories for user %s: %v", user.UserID, err)
		}
	}

	return issueToken(&jwt.UserInfo{
		UserID: user.UserID,
		Name:   user.Name,
	}, c, config)
}

func issueToken(userInfo *jwt.UserInfo, c *fiber.Ctx, config *config.Config) error {
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

func HandleRestore(c *fiber.Ctx, usvc u.Service, config *config.Config) error {
	payload, err := parseGoogleToken(c, config)
	if err != nil {
		return err
	}

	googleID := payload.Subject
	if googleID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "google id not found in token")
	}

	user, err := usvc.RestoreByGoogleID(googleID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to get user")
		}
	}

	return issueToken(&jwt.UserInfo{
		UserID: user.UserID,
		Name:   user.Name,
	}, c, config)
}

func parseGoogleToken(
	c *fiber.Ctx,
	config *config.Config,
) (*idtoken.Payload, error) {
	var req GoogleTokenRequest
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return nil, err
	}

	payload, err := validateIDToken(req.IdToken, config.Google)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid ID token")
	}
	return payload, nil
}
