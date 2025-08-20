package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type StockService interface {
	CreateStock(ctx context.Context, stock dtos.CreateStockReq) (dtos.ID, error)
	UpdateStockFiles(ctx context.Context, stockID dtos.ID, images []string, logo string) error
	GetStocks(ctx context.Context, limit, page int64, search, status string) (dtos.StocksResult, error)
	GetStockByID(ctx context.Context, stockID int64) (dtos.Stock, error)
	UpdateStock(ctx context.Context, stock dtos.UpdateStockReq) (dtos.ID, error)
	DeleteStock(ctx context.Context, id int64) error
	UpdateStockStatus(ctx context.Context, stock dtos.UpdateStockStatus) (dtos.ID, error)
}
