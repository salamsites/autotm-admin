package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type BrandService interface {
	// Body Type
	CreateBodyType(ctx context.Context, bodyType dtos.CreateBodyTypeReq) (dtos.ID, error)
	GetBodyType(ctx context.Context, limit, page int64, category, search string) (dtos.BodyTypeResult, error)
	UpdateBodyType(ctx context.Context, bodyType dtos.UpdateBodyTypeReq) (dtos.ID, error)
	DeleteBodyType(ctx context.Context, id int64) error

	// Brand
	CreateBrand(ctx context.Context, brand dtos.CreateBrandReq) (dtos.ID, error)
	GetBrands(ctx context.Context, limit, page int64, category, search string) (dtos.BrandResult, error)
	UpdateBrand(ctx context.Context, brand dtos.UpdateBrandReq) (dtos.ID, error)
	DeleteBrandCategory(ctx context.Context, id int64, category string) error

	// Model
	CreateModel(ctx context.Context, model dtos.CreateModelReq) (dtos.ID, error)
	GetModels(ctx context.Context, limit, page int64, category, search string) (dtos.ModelResult, error)
	UpdateModel(ctx context.Context, model dtos.UpdateModelReq) (dtos.ID, error)
	DeleteModel(ctx context.Context, id int64) error

	//Description
	CreateDescription(ctx context.Context, description dtos.CreateDescription) (dtos.ID, error)
	GetDescriptions(ctx context.Context, limit, page int64, search, category string) (dtos.DescriptionResult, error)
	UpdateDescription(ctx context.Context, description dtos.UpdateDescription) (dtos.ID, error)
	DeleteDescription(ctx context.Context, id int64) error
}
