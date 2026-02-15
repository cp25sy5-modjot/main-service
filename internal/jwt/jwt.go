package jwt

import (
	"errors"
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	"github.com/cp25sy5-modjot/main-service/internal/shared/config"

	userrepo "github.com/cp25sy5-modjot/main-service/internal/user/repository"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected creates and returns the JWT middleware.
func Protected(secret string, userRepo *userrepo.Repository) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key:    []byte(secret),
			JWTAlg: "HS256",
		},
		Claims:       &Claims{},
		ErrorHandler: jwtErrorHandler,

		SuccessHandler: func(c *fiber.Ctx) error {

			userID, err := GetUserIDFromClaims(c)
			if err != nil {
				return err
			}

			user, err := userRepo.FindByID(userID)
			if err != nil || user == nil {
				return fiber.NewError(fiber.StatusUnauthorized, "User no longer exists")
			}
			if user.Status == e.UserStatusInactive {
				return fiber.NewError(fiber.StatusForbidden, "Account is deactivated")
			}

			c.Locals("authUser", user)

			return c.Next()
		},
	})
}

func jwtErrorHandler(c *fiber.Ctx, err error) error {
	if errors.Is(err, jwt.ErrTokenExpired) {
		return fiber.NewError(fiber.StatusUnauthorized, "Token has expired")
	}
	return fiber.NewError(fiber.StatusUnauthorized, "Invalid or malformed token")
}

// GenerateTokens creates and returns new access and refresh tokens.
func GenerateTokens(user *UserInfo, conf *config.Auth) (accessToken string, refreshToken string, err error) {
	// --- Create Access Token ---
	accessClaims := &Claims{
		Type: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserID,
			Issuer:    conf.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(parseTime(conf.AccessTokenTTL))),
		},
	}
	accessToken, err = createToken(accessClaims, conf.AccessTokenSecret)
	if err != nil {
		return "", "", err
	}

	// --- Create Refresh Token ---
	refreshClaims := &Claims{
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserID,
			Issuer:    conf.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(parseTime(conf.RefreshTokenTTL))),
		},
	}
	refreshToken, err = createToken(refreshClaims, conf.RefreshTokenSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// createToken is a helper function to sign a token with a given secret.
func createToken(claims *Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func parseTime(timeStr string) time.Duration {
	duration, err := time.ParseDuration(timeStr)
	if err != nil {
		panic("Invalid JWT TTL config: " + timeStr)
	}
	return duration
}

func GetUserIDFromClaims(c *fiber.Ctx) (string, error) {

	userVal := c.Locals("user")
	if userVal == nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Missing token")
	}

	token, ok := userVal.(*jwt.Token)
	if !ok {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || claims.Subject == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	return claims.Subject, nil
}
