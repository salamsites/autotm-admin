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

func (s *BrandService) CreateBodyType(ctx context.Context, bodyType dtos.CreateBodyTypeReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(bodyType); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	newBodyType := models.BodyType{
		Name:      bodyType.Name,
		ImagePath: bodyType.ImagePath,
	}

	bodyTypeID, err := s.repo.CreateBodyType(ctx, newBodyType)
	if err != nil {
		s.logger.Errorf("create err: %v", err)
		return bodyTypeID, err
	}
	return bodyTypeID, nil
}

func (s *BrandService) GetBodyType(ctx context.Context, limit, page int64, search string) (dtos.BodyTypeResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	bodyTypes, count, err := s.repo.GetBodyType(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get brands err: %v", err)
		return dtos.BodyTypeResult{}, err
	}
	var dtoBodyTypes []dtos.BodyType
	for _, b := range bodyTypes {
		dtoBodyTypes = append(dtoBodyTypes, dtos.BodyType{
			ID:        b.ID,
			Name:      b.Name,
			ImagePath: b.ImagePath,
		})
	}

	result := dtos.BodyTypeResult{
		BodyTypes: dtoBodyTypes,
		Count:     count,
	}
	return result, nil
}

func (s *BrandService) UpdateBodyType(ctx context.Context, bodyType dtos.UpdateBodyTypeReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(bodyType); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	oldBodyType, err := s.repo.GetBodyTypeByID(ctx, bodyType.ID)
	if err != nil {
		s.logger.Errorf("get old body types err: %v", err)
		return 0, err
	}

	if oldBodyType.ImagePath != bodyType.ImagePath && oldBodyType.ImagePath != "" {
		if err = helpers.DeleteImage(oldBodyType.ImagePath); err != nil {
			s.logger.Errorf("delete old image path err: %v", err)
		}
	}

	newBodyType := models.BodyType{
		ID:        bodyType.ID,
		Name:      bodyType.Name,
		ImagePath: bodyType.ImagePath,
	}

	bodyTypeID, err := s.repo.UpdateBodyType(ctx, newBodyType)
	if err != nil {
		s.logger.Errorf("update body types err: %v", err)
		return bodyTypeID, err
	}
	return bodyTypeID, nil
}

func (s *BrandService) DeleteBodyType(ctx context.Context, id int64) error {
	oldBodyType, err := s.repo.GetBodyTypeByID(ctx, id)
	if err != nil {
		s.logger.Errorf("get old body type err: %v", err)
		return err
	}

	if oldBodyType.ImagePath != "" {
		if err = helpers.DeleteImage(oldBodyType.ImagePath); err != nil {
			s.logger.Errorf("delete old image path err: %v", err)
		}
	}

	deleteID := models.ID{
		ID: id,
	}

	err = s.repo.DeleteBodyType(ctx, deleteID)
	if err != nil {
		s.logger.Errorf("delete body type err: %v", err)
		return err
	}
	return nil
}

func (s *BrandService) CreateBrand(ctx context.Context, brand dtos.CreateBrandReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(brand); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	newBrand := models.Brand{
		Name:       brand.Name,
		LogoPath:   brand.LogoPath,
		Categories: brand.Categories,
	}

	brandID, err := s.repo.CreateBrand(ctx, newBrand)
	if err != nil {
		s.logger.Errorf("create err: %v", err)
		return brandID, err
	}
	return brandID, nil
}

func (s *BrandService) GetBrandsByCategory(ctx context.Context, limit, page int64, categoryType, search string) (dtos.BrandResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	brands, count, err := s.repo.GetBrandsByCategory(ctx, limit, offset, categoryType, search)
	if err != nil {
		s.logger.Errorf("get brands err: %v", err)
		return dtos.BrandResult{}, err
	}
	var dtoBrands []dtos.Brand
	for _, b := range brands {
		dtoBrands = append(dtoBrands, dtos.Brand{
			ID:         b.ID,
			Name:       b.Name,
			LogoPath:   b.LogoPath,
			Categories: b.Categories,
		})
	}

	result := dtos.BrandResult{
		Brands: dtoBrands,
		Count:  count,
	}
	return result, nil
}

func (s *BrandService) UpdateBrand(ctx context.Context, brand dtos.UpdateBrandReq) (int64, error) {
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

func (s *BrandService) DeleteBrandCategory(ctx context.Context, id int64, category string) error {
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
		ID:       id,
		Category: category,
	}

	err = s.repo.DeleteBrandCategory(ctx, deleteID)
	if err != nil {
		s.logger.Errorf("delete brand err: %v", err)
		return err
	}
	return nil
}

func (s *BrandService) CreateModel(ctx context.Context, model dtos.CreateModelReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(model); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	newModel := models.Model{
		Name:       model.Name,
		BrandID:    model.BrandID,
		BodyTypeID: model.BodyTypeID,
	}

	modelID, err := s.repo.CreateModel(ctx, newModel)
	if err != nil {
		s.logger.Errorf("create model err: %v", err)
		return modelID, err
	}
	return modelID, nil
}

func (s *BrandService) GetModels(ctx context.Context, limit, page int64, search string) (dtos.ModelResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	brandModels, count, err := s.repo.GetModels(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get models err: %v", err)
		return dtos.ModelResult{}, err
	}
	var dtoModels []dtos.Model
	for _, b := range brandModels {
		dtoModels = append(dtoModels, dtos.Model{
			ID:           b.ID,
			Name:         b.Name,
			LogoPath:     b.LogoPath,
			BrandID:      b.BrandID,
			BrandName:    b.BrandName,
			BodyTypeID:   b.BodyTypeID,
			BodyTypeName: b.BodyTypeName,
		})
	}

	result := dtos.ModelResult{
		Models: dtoModels,
		Count:  count,
	}
	return result, nil
}

func (s *BrandService) UpdateModel(ctx context.Context, model dtos.UpdateModelReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(model); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	newModel := models.Model{
		ID:         model.ID,
		Name:       model.Name,
		BrandID:    model.BrandID,
		BodyTypeID: model.BodyTypeID,
	}

	brandID, err := s.repo.UpdateModel(ctx, newModel)
	if err != nil {
		s.logger.Errorf("update model err: %v", err)
		return brandID, err
	}
	return brandID, nil
}

func (s *BrandService) DeleteModel(ctx context.Context, id int64) error {
	deleteID := models.ID{
		ID: id,
	}

	err := s.repo.DeleteModel(ctx, deleteID)
	if err != nil {
		s.logger.Errorf("delete model err: %v", err)
		return err
	}
	return nil
}
