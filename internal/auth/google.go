package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/cp25sy5-modjot/main-service/internal/config"

	r "github.com/cp25sy5-modjot/main-service/internal/response"
	u "github.com/cp25sy5-modjot/main-service/internal/user"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/idtoken"
)

func HandleGoogleTokenExchange(c *fiber.Ctx, service *u.Service, config *config.Config) error {
	var reqData GoogleTokenRequest
	if err := c.BodyParser(&reqData); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return r.BadRequest(c, "Cannot parse request body", err)
	}

	idToken, errMsg := exchangeToken(reqData, config)
	if errMsg != "" {
		return r.InternalServerError(c, errMsg)
	}

	payload, err := validateIDToken(idToken, config.Google)
	if err != nil {
		log.Printf("Error validating ID token: %v", err)
		return r.Unauthorized(c, "Invalid ID token")
	}

	userInfo := getUserInfoFromPayload(payload, service)

	accessToken, refreshToken, err := GenerateTokens(userInfo, config.Auth)
	if err != nil {
		return r.InternalServerError(c, "Failed to generate tokens")
	}

	return r.OK(c, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, "Login successful")
}

func requestBuilder(reqData GoogleTokenRequest, config *config.Google) url.Values {
	formData := url.Values{}
	formData.Set("code", reqData.Code)
	formData.Set("code_verifier", reqData.CodeVerifier)
	formData.Set("client_id", config.ClientID)
	formData.Set("redirect_uri", config.RedirectURL)
	formData.Set("grant_type", "authorization_code")
	return formData
}

func exchangeToken(reqData GoogleTokenRequest, config *config.Config) (string, string) {
	formData := requestBuilder(reqData, config.Google)

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", formData)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", "Failed to exchange token with Google"
	}
	defer resp.Body.Close()

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", "failed to parse Google's response"
	}

	idTokenString, ok := tokenResponse["id_token"].(string)
	if !ok {
		return "", "id_token not found in response"
	}
	return idTokenString, ""
}

func validateIDToken(idToken string, config *config.Google) (*idtoken.Payload, error) {
	ctx := context.Background()
	payload, err := idtoken.Validate(ctx, idToken, config.ClientID)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func getUserInfoFromPayload(payload *idtoken.Payload, service *u.Service) *UserInfo {
	userID := payload.Subject

	user, err := service.GetByID(userID)
	if err != nil {
		name := payload.Claims["given_name"].(string)
		if name == "" {
			name = payload.Claims["name"].(string)
		}
		service.Create(&u.UserInsertReq{
			UserID: userID,
			Email:  payload.Claims["email"].(string),
			Name:   name,
		})
		log.Printf("Created new user!")
		user, _ = service.GetByID(userID)
	}

	return &UserInfo{
		UserID: user.UserID,
		Name:   user.Name,
	}
}
