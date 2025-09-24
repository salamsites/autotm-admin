package storage

import (
	"autotm-admin/internal/models"
	"context"

	"github.com/salamsites/minio-pkg/util"
)

type StockRepository interface {
	CreateStock(ctx context.Context, stock models.Stock) (int64, error)
	UpdateStockImages(ctx context.Context, stockID int64, images util.Media) error
	UpdateStockLogo(ctx context.Context, stockID int64, logo interface{}) error
	GetStocks(ctx context.Context, limit, page int64, search, status string) ([]models.Stock, int64, error)
	GetStockByID(ctx context.Context, stockID int64) (models.Stock, error)
	UpdateStock(ctx context.Context, autoStore models.Stock) (int64, error)
	DeleteStock(ctx context.Context, id int64) error
	UpdateStockStatus(ctx context.Context, id int64, status string) (int64, error)
	GetStockFollowers(ctx context.Context, stockId int64) ([]int64, error)
	DeleteStockLogo(ctx context.Context, stockId int64) error
	DeleteStockImage(ctx context.Context, stockId int64, path string) error
}
