package jwt

import (
	"time"

	"github.com/cp25sy5-modjot/main-service/internal/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected creates and returns the JWT middleware.
func Protected(AccessTokenSecret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(AccessTokenSecret), JWTAlg: "HS256"},
		Claims:       &Claims{},
		ErrorHandler: jwtErrorHandler,
	})
}

func jwtErrorHandler(c *fiber.Ctx, err error) error {
	// Check for a specific error type, like an expired token.
	if err.Error() == "token is expired" {
		return fiber.NewError(fiber.StatusUnauthorized, "Token has expired")
	}

	// For all other errors (missing, malformed, invalid signature)
	return fiber.NewError(fiber.StatusUnauthorized, "Invalid or malformed token")
}

// GenerateTokens creates and returns new access and refresh tokens.
func GenerateTokens(user *UserInfo, conf *config.Auth) (accessToken string, refreshToken string, err error) {
	// --- Create Access Token ---
	accessClaims := &Claims{
		Name: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserID,
			Issuer:    conf.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(parseTime(conf.AccessTokenTTL))),
		},
	}
	accessToken, err = createToken(accessClaims, conf.AccessTokenSecret)
	if err != nil {
		return "", "", err
	}

	// --- Create Refresh Token ---
	refreshClaims := &Claims{
		Name: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserID,
			Issuer:    conf.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(parseTime(conf.RefreshTokenTTL))),
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
		return 0
	}
	return duration
}

func GetUserIDFromClaims(c *fiber.Ctx) (string, error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(*Claims)
	if claims == nil || claims.Subject == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Failed to get user ID from claims")
	}
	return claims.Subject, nil
}
