package services

import (
	"autotm-admin/internal/configs"
	"autotm-admin/internal/dtos"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	slog "github.com/salamsites/package-log"
	"io"
	"net/http"
	"net/url"
	"time"
)

type UserService struct {
	logger *slog.Logger
	cfg    *configs.Config
}

func NewUserService(cfg *configs.Config, logger *slog.Logger) *UserService {
	return &UserService{
		logger: logger,
		cfg:    cfg,
	}
}

func (s *UserService) GetUsers(ctx context.Context, limit, page int64, search string) ([]dtos.GetUsers, int64, error) {
	searchParam := url.QueryEscape(search)

	urlUser := fmt.Sprintf("%s/users/get-users?limit=%d&page=%d&search=%s", s.cfg.UserServiceURL, limit, page, searchParam)
	s.logger.Infof("Calling URL: %s", urlUser)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlUser, nil)
	if err != nil {
		s.logger.Errorf("failed to build request: %v", err)
		return nil, 0, err
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		s.logger.Errorf("http request failed: %v", err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Errorf("failed to read response body: %v", err)
		return nil, 0, err
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
		s.logger.Error(errMsg)
		return nil, 0, fmt.Errorf(errMsg)
	}

	var responseWrapper struct {
		Data dtos.GetUserResult `json:"Data"`
	}
	if err = json.Unmarshal(bodyBytes, &responseWrapper); err != nil {
		s.logger.Errorf("failed to decode response body: %v", err)
		return nil, 0, err
	}

	users := responseWrapper.Data.Users
	count := responseWrapper.Data.Count

	return users, count, nil
}

func (s *UserService) GetUserByIds(ctx context.Context, ids dtos.GetUserByIDsReq) ([]dtos.GetUsers, error) {
	urlUser := fmt.Sprintf("%s/users/get-by-ids", s.cfg.UserServiceURL)

	body, err := json.Marshal(ids)
	if err != nil {
		s.logger.Errorf("failed to build request: %s", err)
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", urlUser, bytes.NewBuffer(body))
	if err != nil {
		s.logger.Errorf("failed to create request: %s", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Errorf("http request failed: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Errorf("unexpected status: %d", resp.StatusCode)
		return nil, fmt.Errorf("users not found")
	}

	var result struct {
		Data []dtos.GetUsers `json:"Data"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		s.logger.Errorf("failed to decode response: %s", err)
		return nil, err
	}

	users := result.Data
	return users, nil
}
