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
	var req GoogleTokenRequest
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	payload, err := validateIDToken(req.IdToken, config.Google)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid ID token")
	}

	userInfo, err := getUserInfoFromPayload(payload, usvc, csvc)
	if err != nil {
		return err
	}
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

func getUserInfoFromPayload(payload *idtoken.Payload, usvc u.Service, csvc c.Service) (*jwt.UserInfo, error) {
	if payload == nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid token payload")
	}

	googleID := payload.Subject
	if googleID == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "google id not found in token")
	}

	// 1. ลองหา user จาก GoogleID
	user, err := usvc.GetByGoogleID(googleID)
	if err != nil {
		// ถ้าเป็น error อื่นที่ไม่ใช่ not found → คืน error เลย
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to get user")
		}
	}

	// 2. ถ้ายังไม่มี user (nil) → สร้างใหม่
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
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create user")
		}
		log.Printf("created new user with googleID %s", googleID)

		// Create default categories for the new user
		if err := csvc.CreateDefaultCategories(user.UserID); err != nil {
			log.Printf("failed to create default categories for user %s: %v", user.UserID, err)
		}
	}

	// 3. ตรงนี้มั่นใจได้แล้วว่า user != nil
	return &jwt.UserInfo{
		UserID: user.UserID,
		Name:   user.Name,
	}, nil
}
