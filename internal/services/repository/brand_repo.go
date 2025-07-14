package repository

import (
	"autotm-admin/internal/dtos"
	"context"
	"mime/multipart"
)

type BrandService interface {
	UploadImage(file multipart.File, header *multipart.FileHeader) (string, error)
	CreateBrand(ctx context.Context, brand dtos.CreateBrandReq) (int64, error)
	GetBrands(ctx context.Context, limit, page int64, search string) (dtos.BrandResult, error)
	UpdateBrand(ctx context.Context, brand dtos.Brand) (int64, error)
	DeleteBrand(ctx context.Context, id int64) error
	CreateBrandModel(ctx context.Context, model dtos.CreateBrandModelReq) (int64, error)
	GetBrandModels(ctx context.Context, limit, page int64, search string) (dtos.BrandModelResult, error)
	UpdateBrandModel(ctx context.Context, model dtos.BrandModel) (int64, error)
	DeleteBrandModel(ctx context.Context, id int64) error
}
