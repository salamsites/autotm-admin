package services

import (
	"autotm-admin/internal/configs"
	"autotm-admin/internal/dtos"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	slog "github.com/salamsites/package-log"
)

type PushService struct {
	logger *slog.Logger
	cfg    *configs.Config
	client *http.Client
}

func NewPushService(logger *slog.Logger, cfg *configs.Config) *PushService {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	return &PushService{
		logger: logger,
		cfg:    cfg,
		client: client,
	}
}

func (s *PushService) SendMultiPush(ctx context.Context, pushData dtos.ReqSendPushDTO) error {
	url := fmt.Sprintf("%s/push/send-multi-push", s.cfg.PushService)

	payload, err := json.Marshal(pushData)
	if err != nil {
		return fmt.Errorf("marshal payload error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("http post error: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body error: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("push service error: %s", string(body))
	}

	return nil
}
