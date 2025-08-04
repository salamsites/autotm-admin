package storage

import (
	"autotm-admin/internal/models"
	"context"
)

type AutoStoreRepository interface {
	CreateAutoStore(ctx context.Context, autoStore models.AutoStore) (int64, error)
	GetAutoStores(ctx context.Context, limit, page int64, search string) ([]models.AutoStore, int64, error)
	UpdateAutoStore(ctx context.Context, autoStore models.AutoStore) (int64, error)
	DeleteAutoStore(ctx context.Context, id models.ID) error
}
