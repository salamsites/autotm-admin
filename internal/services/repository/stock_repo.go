package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type StockService interface {
	CreateStock(ctx context.Context, stock dtos.CreateStockReq) (int64, error)
	UpdateStockFiles(ctx context.Context, stockID int64, images []string, logo string) error
	GetStocks(ctx context.Context, limit, page int64, search string) (dtos.StocksResult, error)
	UpdateStock(ctx context.Context, stock dtos.UpdateStockReq) (dtos.ID, error)
	DeleteStock(ctx context.Context, id int64) error
}
