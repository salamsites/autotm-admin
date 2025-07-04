package repository

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/models"
	"context"
	"mime/multipart"
)

type BrandService interface {
	UploadImage(file multipart.File, header *multipart.FileHeader) (string, error)
	CreateBrand(ctx context.Context, brand dtos.V1BrandDTO) (int64, error)
	GetBrands(ctx context.Context, limit, page int64, search string) (models.BrandResult, error)
	UpdateBrand(ctx context.Context, brand dtos.V1BrandDTO) (int64, error)
	DeleteBrand(ctx context.Context, id int64) error
	CreateBrandModel(ctx context.Context, model dtos.V1BrandModelDTO) (int64, error)
	GetBrandModels(ctx context.Context, limit, page int64, search string) (models.BrandModelResult, error)
	UpdateBrandModel(ctx context.Context, model dtos.V1BrandModelDTO) (int64, error)
	DeleteBrandModel(ctx context.Context, id int64) error
}
