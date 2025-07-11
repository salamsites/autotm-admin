package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type RegionsService interface {
	CreateRegion(ctx context.Context, region dtos.Region) (int64, error)
	GetAllRegions(ctx context.Context, limit, page int64, search string) (dtos.RegionResult, error)
	UpdateRegion(ctx context.Context, region dtos.Region) (int64, error)
	DeleteRegion(ctx context.Context, id int64) error

	// Cities
	CreateCity(ctx context.Context, city dtos.City) (int64, error)
	GetAllCities(ctx context.Context, limit, page int64, search string) (dtos.CityResult, error)
	UpdateCity(ctx context.Context, region dtos.City) (int64, error)
	DeleteCity(ctx context.Context, id int64) error
}
