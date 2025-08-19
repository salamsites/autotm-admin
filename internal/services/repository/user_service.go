package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type UserService interface {
	GetUsersFromUserService(ctx context.Context, limit, page int64, search string) (dtos.GetUsersResult, error)
}
