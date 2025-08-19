package storage

import (
	"autotm-admin/internal/models"
	"context"
)

type StockRepository interface {
	CreateStock(ctx context.Context, stock models.Stock) (int64, error)
	UpdateStockImages(ctx context.Context, stockID int64, images []string) error
	UpdateStockLogo(ctx context.Context, stockID int64, logo string) error
	GetStocks(ctx context.Context, limit, page int64, search string) ([]models.Stock, int64, error)
	GetStockByID(ctx context.Context, stockID int64) (models.Stock, error)
	UpdateStock(ctx context.Context, autoStore models.Stock) (int64, error)
	DeleteStock(ctx context.Context, id int64) error
	UpdateStockStatus(ctx context.Context, id int64, status string) (int64, error)
}
