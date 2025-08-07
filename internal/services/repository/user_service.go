package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type UserService interface {
	GetUsers(ctx context.Context, limit, page int64, search string) ([]dtos.GetUsers, int64, error)
	GetUserByIds(ctx context.Context, ids dtos.GetUserByIDsReq) ([]dtos.GetUsers, error)
}
