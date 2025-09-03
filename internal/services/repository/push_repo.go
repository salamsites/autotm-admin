package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type PushService interface {
	SendMultiPush(ctx context.Context, req dtos.ReqSendPushDTO) error
}
