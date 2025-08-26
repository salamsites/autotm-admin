package repository

import "autotm-admin/internal/dtos"

type PushService interface {
	SendPush(req dtos.ReqSendPushDTO) error
}
