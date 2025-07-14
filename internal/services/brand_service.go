package services

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"autotm-admin/internal/models"
	"autotm-admin/internal/repository/storage"
	"context"
	slog "github.com/salamsites/package-log"
	"mime/multipart"
)

type BrandService struct {
	logger *slog.Logger
	repo   storage.BrandRepository
}

func NewBrandService(logger *slog.Logger, repo storage.BrandRepository) *BrandService {
	return &BrandService{
		logger: logger,
		repo:   repo,
	}
}

func (s *BrandService) UploadImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	imagePath, err := helpers.UploadImage(file, header)
	if err != nil {
		s.logger.Errorf("upload image err: %v", err)
		return imagePath, err
	}
	return imagePath, nil
}

func (s *BrandService) CreateBrand(ctx context.Context, brand dtos.CreateBrandReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(brand); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	newBrand := models.Brand{
		Name:     brand.Name,
		LogoPath: brand.LogoPath,
	}

	brandID, err := s.repo.CreateBrand(ctx, newBrand)
	if err != nil {
		s.logger.Errorf("create err: %v", err)
		return brandID, err
	}
	return brandID, nil
}

func (s *BrandService) GetBrands(ctx context.Context, limit, page int64, search string) (dtos.BrandResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	brands, count, err := s.repo.GetBrands(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get brands err: %v", err)
		return dtos.BrandResult{}, err
	}
	var dtoBrands []dtos.Brand
	for _, b := range brands {
		dtoBrands = append(dtoBrands, dtos.Brand{
			ID:       b.ID,
			Name:     b.Name,
			LogoPath: b.LogoPath,
		})
	}

	result := dtos.BrandResult{
		Brands: dtoBrands,
		Count:  count,
	}
	return result, nil
}

func (s *BrandService) UpdateBrand(ctx context.Context, brand dtos.Brand) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(brand); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	oldBrand, err := s.repo.GetBrandByID(ctx, brand.ID)
	if err != nil {
		s.logger.Errorf("get old brand err: %v", err)
		return 0, err
	}

	if oldBrand.LogoPath != brand.LogoPath && oldBrand.LogoPath != "" {
		if err := helpers.DeleteImage(oldBrand.LogoPath); err != nil {
			s.logger.Errorf("delete old logo path err: %v", err)
		}
	}

	newBrand := models.Brand{
		ID:       brand.ID,
		Name:     brand.Name,
		LogoPath: brand.LogoPath,
	}

	brandID, err := s.repo.UpdateBrand(ctx, newBrand)
	if err != nil {
		s.logger.Errorf("update brand err: %v", err)
		return brandID, err
	}
	return brandID, nil
}

func (s *BrandService) DeleteBrand(ctx context.Context, id int64) error {
	oldBrand, err := s.repo.GetBrandByID(ctx, id)
	if err != nil {
		s.logger.Errorf("get old brand err: %v", err)
		return err
	}

	if oldBrand.LogoPath != "" {
		if err := helpers.DeleteImage(oldBrand.LogoPath); err != nil {
			s.logger.Errorf("delete old logo path err: %v", err)
		}
	}

	deleteID := models.ID{
		ID: id,
	}

	err = s.repo.DeleteBrand(ctx, deleteID)
	if err != nil {
		s.logger.Errorf("delete brand err: %v", err)
		return err
	}
	return nil
}

func (s *BrandService) CreateBrandModel(ctx context.Context, model dtos.CreateBrandModelReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(model); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	newModel := models.BrandModel{
		Name:    model.Name,
		BrandID: model.BrandID,
	}

	modelID, err := s.repo.CreateBrandModel(ctx, newModel)
	if err != nil {
		s.logger.Errorf("create model err: %v", err)
		return modelID, err
	}
	return modelID, nil
}

func (s *BrandService) GetBrandModels(ctx context.Context, limit, page int64, search string) (dtos.BrandModelResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	brandModels, count, err := s.repo.GetBrandModels(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get brand models err: %v", err)
		return dtos.BrandModelResult{}, err
	}
	var dtoBrands []dtos.BrandModel
	for _, b := range brandModels {
		dtoBrands = append(dtoBrands, dtos.BrandModel{
			ID:        b.ID,
			Name:      b.Name,
			LogoPath:  b.LogoPath,
			BrandID:   b.BrandID,
			BrandName: b.BrandName,
		})
	}

	result := dtos.BrandModelResult{
		BrandModels: dtoBrands,
		Count:       count,
	}
	return result, nil
}

func (s *BrandService) UpdateBrandModel(ctx context.Context, model dtos.BrandModel) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(model); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	newModel := models.BrandModel{
		ID:       model.ID,
		Name:     model.Name,
		LogoPath: model.LogoPath,
		BrandID:  model.BrandID,
	}

	brandID, err := s.repo.UpdateBrandModel(ctx, newModel)
	if err != nil {
		s.logger.Errorf("update brand model err: %v", err)
		return brandID, err
	}
	return brandID, nil
}

func (s *BrandService) DeleteBrandModel(ctx context.Context, id int64) error {
	deleteID := models.ID{
		ID: id,
	}

	err := s.repo.DeleteBrandModel(ctx, deleteID)
	if err != nil {
		s.logger.Errorf("delete brand model err: %v", err)
		return err
	}
	return nil
}
