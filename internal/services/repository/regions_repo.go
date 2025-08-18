package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type RegionsService interface {
	// Regions
	CreateRegion(ctx context.Context, region dtos.CreateRegionReq) (dtos.ID, error)
	GetAllRegions(ctx context.Context, limit, page int64, search string) (dtos.RegionResult, error)
	UpdateRegion(ctx context.Context, region dtos.UpdateRegionReq) (dtos.ID, error)
	DeleteRegion(ctx context.Context, id int64) error

	// Cities
	CreateCity(ctx context.Context, city dtos.CreateCityReq) (dtos.ID, error)
	GetAllCities(ctx context.Context, limit, page int64, search string, regionIds []int64) (dtos.CityResult, error)
	UpdateCity(ctx context.Context, region dtos.UpdateCityReq) (dtos.ID, error)
	DeleteCity(ctx context.Context, id int64) error
}
