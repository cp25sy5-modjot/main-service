package pushsvc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	repo "github.com/cp25sy5-modjot/main-service/internal/push_tokens/repository"
)

type Service interface {
	Register(ctx context.Context, userID, token, platform string) error
	Send(ctx context.Context, userID, title, body string) error
}

type service struct {
	repo *repo.Repository
}

func NewService(r *repo.Repository) Service {
	return &service{repo: r}
}

func (s *service) Register(ctx context.Context, userID, token, platform string) error {
	return s.repo.Save(ctx, userID, token, platform)
}

func (s *service) Send(ctx context.Context, userID, title, body string) error {
	tokens, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		return nil
	}

	var messages []map[string]interface{}
	for _, t := range tokens {
		messages = append(messages, map[string]interface{}{
			"to":    t,
			"title": title,
			"body":  body,
		})
	}

	payload, _ := json.Marshal(messages)

	resp, err := http.Post(
		"https://exp.host/--/api/v2/push/send",
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			Status  string `json:"status"`
			Details struct {
				Error string `json:"error"`
			} `json:"details"`
		} `json:"data"`
	}

	json.NewDecoder(resp.Body).Decode(&result)

	// cleanup invalid tokens
	for i, r := range result.Data {
		if r.Status == "error" && r.Details.Error == "DeviceNotRegistered" {
			_ = s.repo.DeleteByToken(ctx, tokens[i])
		}
	}

	return nil
}