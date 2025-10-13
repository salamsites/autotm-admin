package storage

import (
	"autotm-admin/internal/models"
	"context"
)

type BrandRepository interface {
	//Body Type
	CreateBodyType(ctx context.Context, bodyType models.BodyType) (int64, error)
	GetBodyType(ctx context.Context, limit, page int64, category, search string) ([]models.BodyType, error)
	GetBodyTypeCount(ctx context.Context, category, search string) (int64, error)
	UpdateBodyType(ctx context.Context, bodyType models.BodyType) (int64, error)
	GetBodyTypeByID(ctx context.Context, id int64) (models.BodyType, error)
	DeleteBodyType(ctx context.Context, id int64) error

	//Brand
	CreateBrand(ctx context.Context, brand models.Brand) (int64, error)
	GetBrands(ctx context.Context, limit, page int64, category, search string) ([]models.Brand, error)
	GetBrandsCount(ctx context.Context, category, search string) (int64, error)
	UpdateBrand(ctx context.Context, brand models.Brand) (int64, error)
	GetBrandByID(ctx context.Context, id int64) (models.Brand, error)
	DeleteBrandCategory(ctx context.Context, id int64, category string) error

	// Model
	CreateModel(ctx context.Context, model models.Model) (int64, error)
	GetModels(ctx context.Context, limit, page int64, category, search string) ([]models.Model, error)
	GetModelsCount(ctx context.Context, category, search string) (int64, error)
	UpdateModel(ctx context.Context, model models.Model) (int64, error)
	DeleteModel(ctx context.Context, id int64) error

	//Description
	CreateDescription(ctx context.Context, description models.Description) (int64, error)
	GetDescriptions(ctx context.Context, limit, page int64, search, category string) ([]models.Description, error)
	GetDescriptionsCount(ctx context.Context, category, search string) (int64, error)
	UpdateDescription(ctx context.Context, description models.Description) (int64, error)
	DeleteDescription(ctx context.Context, id int64, category string) error
}
