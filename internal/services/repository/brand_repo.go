package repository

import (
	"autotm-admin/internal/dtos"
	"context"
	"mime/multipart"
)

type BrandService interface {
	UploadImage(file multipart.File, header *multipart.FileHeader) (string, error)

	// Body Type
	CreateBodyType(ctx context.Context, bodyType dtos.CreateBodyTypeReq) (int64, error)
	GetBodyType(ctx context.Context, limit, page int64, category, search string) (dtos.BodyTypeResult, error)
	UpdateBodyType(ctx context.Context, bodyType dtos.UpdateBodyTypeReq) (int64, error)
	DeleteBodyType(ctx context.Context, id int64) error

	// Brand
	CreateBrand(ctx context.Context, brand dtos.CreateBrandReq) (int64, error)
	GetBrandsByCategory(ctx context.Context, limit, page int64, categoryType, search string) (dtos.BrandResult, error)
	UpdateBrand(ctx context.Context, brand dtos.UpdateBrandReq) (int64, error)
	DeleteBrandCategory(ctx context.Context, id int64, category string) error

	// Model
	CreateModel(ctx context.Context, model dtos.CreateModelReq) (int64, error)
	GetModels(ctx context.Context, limit, page int64, search string) (dtos.ModelResult, error)
	UpdateModel(ctx context.Context, model dtos.UpdateModelReq) (int64, error)
	DeleteModel(ctx context.Context, id int64) error
}
