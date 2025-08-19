package services

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"autotm-admin/internal/models"
	"autotm-admin/internal/repository/storage"
	"context"

	sminio "github.com/salamsites/minio-pkg"
	slog "github.com/salamsites/package-log"
)

type BrandService struct {
	logger          *slog.Logger
	repo            storage.BrandRepository
	minioFileClient sminio.ImageClient
}

func NewBrandService(logger *slog.Logger, repo storage.BrandRepository, minioFileClient sminio.ImageClient) *BrandService {
	return &BrandService{
		logger:          logger,
		repo:            repo,
		minioFileClient: minioFileClient,
	}
}

func (s *BrandService) CreateBodyType(ctx context.Context, bodyType dtos.CreateBodyTypeReq) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(bodyType); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	newBodyType := models.BodyType{
		NameTM:    bodyType.NameTM,
		NameEN:    bodyType.NameEN,
		NameRU:    bodyType.NameRU,
		ImagePath: bodyType.ImagePath,
		Category:  bodyType.Category,
		UploadId:  bodyType.UploadId,
	}

	bodyTypeID, err := s.repo.CreateBodyType(ctx, newBodyType)
	if err != nil {
		s.logger.Errorf("create err: %v", err)
		return id, err
	}

	id.ID = bodyTypeID
	return id, nil
}

func (s *BrandService) GetBodyType(ctx context.Context, limit, page int64, category, search string) (dtos.BodyTypeResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	bodyTypes, count, err := s.repo.GetBodyType(ctx, limit, offset, category, search)
	if err != nil {
		s.logger.Errorf("get brands err: %v", err)
		return dtos.BodyTypeResult{}, err
	}
	var dtoBodyTypes []dtos.BodyType
	for _, b := range bodyTypes {
		dtoBodyTypes = append(dtoBodyTypes, dtos.BodyType{
			ID:        b.ID,
			NameTM:    b.NameTM,
			NameEN:    b.NameEN,
			NameRU:    b.NameRU,
			ImagePath: b.ImagePath,
			UploadId:  b.UploadId,
			Category:  b.Category,
		})
	}

	result := dtos.BodyTypeResult{
		BodyTypes: dtoBodyTypes,
		Count:     count,
	}
	return result, nil
}

func (s *BrandService) UpdateBodyType(ctx context.Context, bodyType dtos.UpdateBodyTypeReq) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(bodyType); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	newBodyType := models.BodyType{
		ID:        bodyType.ID,
		NameTM:    bodyType.NameTM,
		NameEN:    bodyType.NameEN,
		NameRU:    bodyType.NameRU,
		ImagePath: bodyType.ImagePath,
		UploadId:  bodyType.UploadId,
		Category:  bodyType.Category,
	}

	bodyTypeID, err := s.repo.UpdateBodyType(ctx, newBodyType)
	if err != nil {
		s.logger.Errorf("update body types err: %v", err)
		return id, err
	}

	id.ID = bodyTypeID
	return id, nil
}

func (s *BrandService) DeleteBodyType(ctx context.Context, id int64) error {
	err := s.repo.DeleteBodyType(ctx, id)
	if err != nil {
		s.logger.Errorf("delete body type err: %v", err)
		return err
	}
	return nil
}

func (s *BrandService) CreateBrand(ctx context.Context, brand dtos.CreateBrandReq) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(brand); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	newBrand := models.Brand{
		Name:       brand.Name,
		LogoPath:   brand.LogoPath,
		UploadId:   brand.UploadId,
		Categories: brand.Categories,
	}

	brandID, err := s.repo.CreateBrand(ctx, newBrand)
	if err != nil {
		s.logger.Errorf("create err: %v", err)
		return id, err
	}

	id.ID = brandID
	return id, nil
}

func (s *BrandService) GetBrands(ctx context.Context, limit, page int64, category, search string) (dtos.BrandResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	brands, count, err := s.repo.GetBrands(ctx, limit, offset, category, search)
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
			UploadId:   b.UploadId,
			Categories: b.Categories,
		})
	}

	result := dtos.BrandResult{
		Brands: dtoBrands,
		Count:  count,
	}
	return result, nil
}

func (s *BrandService) UpdateBrand(ctx context.Context, brand dtos.UpdateBrandReq) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(brand); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	newBrand := models.Brand{
		ID:         brand.ID,
		Name:       brand.Name,
		LogoPath:   brand.LogoPath,
		UploadId:   brand.UploadId,
		Categories: brand.Categories,
	}

	brandID, err := s.repo.UpdateBrand(ctx, newBrand)
	if err != nil {
		s.logger.Errorf("update brand err: %v", err)
		return id, err
	}

	id.ID = brandID
	return id, nil
}

func (s *BrandService) DeleteBrandCategory(ctx context.Context, id int64, category string) error {
	err := s.repo.DeleteBrandCategory(ctx, id, category)
	if err != nil {
		s.logger.Errorf("delete brand err: %v", err)
		return err
	}
	return nil
}

func (s *BrandService) CreateModel(ctx context.Context, model dtos.CreateModelReq) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(model); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	newModel := models.Model{
		Name:     model.Name,
		BrandID:  model.BrandID,
		Category: model.Category,
	}

	modelID, err := s.repo.CreateModel(ctx, newModel)
	if err != nil {
		s.logger.Errorf("create model err: %v", err)
		return id, err
	}

	id.ID = modelID
	return id, nil
}

func (s *BrandService) GetModels(ctx context.Context, limit, page int64, category, search string) (dtos.ModelResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	brandModels, count, err := s.repo.GetModels(ctx, limit, offset, category, search)
	if err != nil {
		s.logger.Errorf("get models err: %v", err)
		return dtos.ModelResult{}, err
	}
	var dtoModels []dtos.Model
	for _, b := range brandModels {
		dtoModels = append(dtoModels, dtos.Model{
			ID:        b.ID,
			Name:      b.Name,
			LogoPath:  b.LogoPath,
			BrandID:   b.BrandID,
			BrandName: b.BrandName,
			Category:  b.Category,
			UploadId:  b.UploadId,
		})
	}

	result := dtos.ModelResult{
		Models: dtoModels,
		Count:  count,
	}
	return result, nil
}

func (s *BrandService) UpdateModel(ctx context.Context, model dtos.UpdateModelReq) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(model); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	newModel := models.Model{
		ID:       model.ID,
		Name:     model.Name,
		BrandID:  model.BrandID,
		Category: model.Category,
	}

	modelID, err := s.repo.UpdateModel(ctx, newModel)
	if err != nil {
		s.logger.Errorf("update model err: %v", err)
		return id, err
	}
	id.ID = modelID
	return id, nil
}

func (s *BrandService) DeleteModel(ctx context.Context, id int64) error {
	err := s.repo.DeleteModel(ctx, id)
	if err != nil {
		s.logger.Errorf("delete model err: %v", err)
		return err
	}
	return nil
}

func (s *BrandService) CreateDescription(ctx context.Context, description dtos.CreateDescription) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(description); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	newDescription := models.Description{
		NameTM: description.NameTM,
		NameEN: description.NameEN,
		NameRU: description.NameRU,
	}

	descriptionID, err := s.repo.CreateDescription(ctx, newDescription)
	if err != nil {
		s.logger.Errorf("create err: %v", err)
		return id, err
	}

	id.ID = descriptionID
	return id, nil
}

func (s *BrandService) GetDescriptions(ctx context.Context, limit, page int64, search string) (dtos.DescriptionResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	descriptions, count, err := s.repo.GetDescriptions(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get descriptions err: %v", err)
		return dtos.DescriptionResult{}, err
	}
	var dtoDescriptions []dtos.Description
	for _, b := range descriptions {
		dtoDescriptions = append(dtoDescriptions, dtos.Description{
			ID:     b.ID,
			NameTM: b.NameTM,
			NameEN: b.NameEN,
			NameRU: b.NameRU,
		})
	}

	result := dtos.DescriptionResult{
		Descriptions: dtoDescriptions,
		Count:        count,
	}
	return result, nil
}

func (s *BrandService) UpdateDescription(ctx context.Context, description dtos.UpdateDescription) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(description); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	newDescription := models.Description{
		ID:     description.ID,
		NameTM: description.NameTM,
		NameEN: description.NameEN,
		NameRU: description.NameRU,
	}

	descriptionID, err := s.repo.UpdateDescription(ctx, newDescription)
	if err != nil {
		s.logger.Errorf("update descriptions err: %v", err)
		return id, err
	}

	id.ID = descriptionID
	return id, nil
}

func (s *BrandService) DeleteDescription(ctx context.Context, id int64) error {
	err := s.repo.DeleteDescription(ctx, id)
	if err != nil {
		s.logger.Errorf("delete description err: %v", err)
		return err
	}
	return nil
}
