package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type RegionsService interface {
	// Regions
	CreateRegion(ctx context.Context, region dtos.CreateRegionReq) (int64, error)
	GetAllRegions(ctx context.Context, limit, page int64, search string) (dtos.RegionResult, error)
	UpdateRegion(ctx context.Context, region dtos.UpdateRegionReq) (int64, error)
	DeleteRegion(ctx context.Context, id int64) error

	// Cities
	CreateCity(ctx context.Context, city dtos.CreateCityReq) (int64, error)
	GetAllCities(ctx context.Context, limit, page int64, search string) (dtos.CityResult, error)
	UpdateCity(ctx context.Context, region dtos.UpdateCityReq) (int64, error)
	DeleteCity(ctx context.Context, id int64) error
}
