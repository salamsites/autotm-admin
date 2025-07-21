package storage

import (
	"autotm-admin/internal/models"
	"context"
)

type BrandRepository interface {
	//Body Type
	CreateBodyType(ctx context.Context, bodyType models.BodyType) (int64, error)
	GetBodyType(ctx context.Context, limit, page int64, search string) ([]models.BodyType, int64, error)
	UpdateBodyType(ctx context.Context, bodyType models.BodyType) (int64, error)
	GetBodyTypeByID(ctx context.Context, id int64) (models.BodyType, error)
	DeleteBodyType(ctx context.Context, id models.ID) error

	//Brand
	CreateBrand(ctx context.Context, brand models.Brand) (int64, error)
	GetBrandsByCategory(ctx context.Context, limit, page int64, category, search string) ([]models.Brand, int64, error)
	UpdateBrand(ctx context.Context, brand models.Brand) (int64, error)
	GetBrandByID(ctx context.Context, id int64) (models.Brand, error)
	DeleteBrandCategory(ctx context.Context, id models.ID) error

	// Model
	CreateModel(ctx context.Context, model models.Model) (int64, error)
	GetModels(ctx context.Context, limit, page int64, search string) ([]models.Model, int64, error)
	UpdateModel(ctx context.Context, model models.Model) (int64, error)
	DeleteModel(ctx context.Context, id models.ID) error
}
