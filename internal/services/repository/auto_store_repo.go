package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type AutoStoreService interface {
	CreateAutoStore(ctx context.Context, autoStore dtos.CreateAutoStoreReq) (int64, error)
	GetUsersFromUserService(ctx context.Context, limit, page int64, search string) (dtos.GetUserResult, error)
	GetAutoStores(ctx context.Context, limit, page int64, search string) (dtos.AutoStoresResult, error)
	UpdateAutoStore(ctx context.Context, autoStore dtos.UpdateAutoStoreReq) (dtos.ID, error)
	DeleteAutoStore(ctx context.Context, id int64) error
}
