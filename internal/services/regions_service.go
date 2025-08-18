package services

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"autotm-admin/internal/models"
	"autotm-admin/internal/repository/storage"
	"context"

	slog "github.com/salamsites/package-log"
)

type RegionsService struct {
	logger *slog.Logger
	repo   storage.RegionsRepository
}

func NewRegionsService(logger *slog.Logger, repo storage.RegionsRepository) *RegionsService {
	return &RegionsService{
		logger: logger,
		repo:   repo,
	}
}

func (s *RegionsService) CreateRegion(ctx context.Context, region dtos.CreateRegionReq) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(region); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	newRegion := models.Region{
		NameTM: region.NameTM,
		NameEN: region.NameEN,
		NameRU: region.NameRu,
	}

	regionID, err := s.repo.CreateRegion(ctx, newRegion)
	if err != nil {
		s.logger.Errorf("create err: %v", err)
		return id, err
	}
	id.ID = regionID
	return id, nil
}

func (s *RegionsService) GetAllRegions(ctx context.Context, limit, page int64, search string) (dtos.RegionResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	regions, count, err := s.repo.GetAllRegions(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get all regions err: %v", err)
		return dtos.RegionResult{}, err
	}
	var dtoRegions []dtos.Region
	for _, b := range regions {
		dtoRegions = append(dtoRegions, dtos.Region{
			ID:     b.ID,
			NameTM: b.NameTM,
			NameEN: b.NameEN,
			NameRu: b.NameRU,
		})
	}

	result := dtos.RegionResult{
		Regions: dtoRegions,
		Count:   count,
	}
	return result, nil
}

func (s *RegionsService) UpdateRegion(ctx context.Context, region dtos.UpdateRegionReq) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(region); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	newRegion := models.Region{
		ID:     region.ID,
		NameTM: region.NameTM,
		NameEN: region.NameEN,
		NameRU: region.NameRu,
	}

	regionID, err := s.repo.UpdateRegion(ctx, newRegion)
	if err != nil {
		s.logger.Errorf("update region err: %v", err)
		return id, err
	}
	id.ID = regionID
	return id, nil
}

func (s *RegionsService) DeleteRegion(ctx context.Context, id int64) error {
	deleteID := models.ID{
		ID: id,
	}

	err := s.repo.DeleteRegion(ctx, deleteID)
	if err != nil {
		s.logger.Errorf("delete region err: %v", err)
		return err
	}
	return nil
}

// Cities
func (s *RegionsService) CreateCity(ctx context.Context, city dtos.CreateCityReq) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(city); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	newCity := models.City{
		NameTM:   city.NameTM,
		NameEN:   city.NameEN,
		NameRU:   city.NameRu,
		RegionID: city.RegionID,
	}

	cityID, err := s.repo.CreateCity(ctx, newCity)
	if err != nil {
		s.logger.Errorf("create err: %v", err)
		return id, err
	}
	id.ID = cityID
	return id, nil
}

func (s *RegionsService) GetAllCities(ctx context.Context, limit, page int64, search string) (dtos.CityResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	cities, count, err := s.repo.GetAllCities(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get all cities err: %v", err)
		return dtos.CityResult{}, err
	}
	var dtoCities []dtos.City
	for _, b := range cities {
		dtoCities = append(dtoCities, dtos.City{
			ID:           b.ID,
			NameTM:       b.NameTM,
			NameEN:       b.NameEN,
			NameRu:       b.NameRU,
			RegionID:     b.RegionID,
			RegionNameTM: b.RegionNameTM,
			RegionNameEN: b.RegionNameEN,
			RegionNameRU: b.RegionNameRU,
		})
	}

	result := dtos.CityResult{
		Cities: dtoCities,
		Count:  count,
	}
	return result, nil
}

func (s *RegionsService) UpdateCity(ctx context.Context, city dtos.UpdateCityReq) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(city); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	newCity := models.City{
		ID:       city.ID,
		NameTM:   city.NameTM,
		NameEN:   city.NameEN,
		NameRU:   city.NameRu,
		RegionID: city.RegionID,
	}

	cityID, err := s.repo.UpdateCity(ctx, newCity)
	if err != nil {
		s.logger.Errorf("update city err: %v", err)
		return id, err
	}
	id.ID = cityID
	return id, nil
}

func (s *RegionsService) DeleteCity(ctx context.Context, id int64) error {
	deleteID := models.ID{
		ID: id,
	}

	err := s.repo.DeleteCity(ctx, deleteID)
	if err != nil {
		s.logger.Errorf("delete city err: %v", err)
		return err
	}
	return nil
}
