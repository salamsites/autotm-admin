package storage

import (
	"autotm-admin/internal/models"
	"context"
)

type BrandRepository interface {
	CreateBrand(ctx context.Context, brand models.Brand) (int64, error)
	GetBrands(ctx context.Context, limit, page int64, search string) ([]models.Brand, int64, error)
	UpdateBrand(ctx context.Context, brand models.Brand) (int64, error)
	GetBrandByID(ctx context.Context, id int64) (models.Brand, error)
	DeleteBrand(ctx context.Context, id models.ID) error
	CreateBrandModel(ctx context.Context, model models.BrandModel) (int64, error)
	GetBrandModels(ctx context.Context, limit, page int64, search string) ([]models.BrandModel, int64, error)
	UpdateBrandModel(ctx context.Context, model models.BrandModel) (int64, error)
	GetBrandModelByID(ctx context.Context, id int64) (models.BrandModel, error)
	DeleteBrandModel(ctx context.Context, id models.ID) error
}
