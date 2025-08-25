package services

import (
	"autotm-admin/internal/configs"
	"autotm-admin/internal/dtos"
	"bytes"
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
}

func NewPushService(logger *slog.Logger, cfg *configs.Config) *PushService {
	return &PushService{
		logger: logger,
		cfg:    cfg,
	}
}

func (s *PushService) SendPush(req dtos.ReqSendPushDTO) error {
	url := fmt.Sprintf("%s/push/send-push", s.cfg.PushService)

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal payload error: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(payload))
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
