package repository

import (
	"autotm-admin/internal/dtos"
	"context"
	"net/http"
)

type StockService interface {
	CreateStock(ctx context.Context, stock dtos.CreateStockReq, r *http.Request) (dtos.ID, error)
	GetStocks(ctx context.Context, limit, page int64, search, status string) (dtos.StocksResult, error)
	GetStockByID(ctx context.Context, stockID int64) (dtos.Stock, error)
	UpdateStock(ctx context.Context, stock dtos.UpdateStockReq, r *http.Request) (dtos.ID, error)
	DeleteStock(ctx context.Context, id int64) error
	UpdateStockStatus(ctx context.Context, stock dtos.UpdateStockStatus) (dtos.ID, error)
	DeleteStockLogo(ctx context.Context, stockId int64) error
	DeleteStockImage(ctx context.Context, stockId, generateId int64) error
}
