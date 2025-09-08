package elasticsearch

import (
	"autotm-admin/internal/models"
	"context"
	"fmt"

	"github.com/Hajymuhammet/elasticsearch-package/elasticsearch"
	mapping2 "github.com/Hajymuhammet/elasticsearch-package/pkg/mapping"
	esService "github.com/Hajymuhammet/elasticsearch-package/service"
	es "github.com/elastic/go-elasticsearch/v8"
)

type StockESService struct {
	service *esService.Service[models.ESStock]
}

func NewStockESService(client *es.Client) *StockESService {
	repo := elasticsearch.NewRepository[models.ESStock](client)
	service := esService.NewService[models.ESStock](repo)
	return &StockESService{service: service}
}

func (s *StockESService) EnsureIndex(ctx context.Context) error {
	sample := models.ESStock{}
	mapping := mapping2.BuildMappingFromStruct(sample)
	return s.service.EnsureIndex(ctx, "stocks", mapping)
}

func (s *StockESService) IndexStock(ctx context.Context, stock models.ESStock) error {
	return s.service.Index(ctx, "stocks", fmt.Sprintf("%d", stock.ID), stock)
}

func (s *StockESService) BulkIndexStocks(ctx context.Context, stocks []models.ESStock) error {
	return s.service.BulkIndex(ctx, "stocks", stocks, func(stock models.ESStock) string {
		return fmt.Sprintf("%d", stock.ID)
	})
}

func (s *StockESService) SearchStocks(ctx context.Context, query map[string]any, from, size int, sort []map[string]any) ([]models.ESStock, int64, error) {
	return s.service.Search(ctx, "stocks", query, from, size, sort)
}

func (s *StockESService) GetStockByID(ctx context.Context, id int64) (*models.ESStock, error) {
	return s.service.GetByID(ctx, "stocks", fmt.Sprintf("%d", id))
}
