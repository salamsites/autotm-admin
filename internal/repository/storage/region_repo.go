package storage

import (
	"autotm-admin/internal/models"
	"context"
)

type RegionsRepository interface {
	CreateRegion(ctx context.Context, model models.Region) (int64, error)
	GetAllRegions(ctx context.Context, limit, page int64, search string) ([]models.Region, int64, error)
	UpdateRegion(ctx context.Context, region models.Region) (int64, error)
	DeleteRegion(ctx context.Context, id models.ID) error

	//Cities
	CreateCity(ctx context.Context, model models.City) (int64, error)
	GetAllCities(ctx context.Context, limit, page int64, search string, regionIds []int64) ([]models.City, int64, error)
	UpdateCity(ctx context.Context, region models.City) (int64, error)
	DeleteCity(ctx context.Context, id models.ID) error
}
