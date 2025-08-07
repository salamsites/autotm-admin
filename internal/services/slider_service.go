package services

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"autotm-admin/internal/models"
	"autotm-admin/internal/repository/storage"
	"context"
	slog "github.com/salamsites/package-log"
)

type SlidersService struct {
	logger *slog.Logger
	repo   storage.SlidersRepository
}

func NewSlidersService(logger *slog.Logger, repo storage.SlidersRepository) *SlidersService {
	return &SlidersService{
		logger: logger,
		repo:   repo,
	}
}

func (s *SlidersService) CreateSlider(ctx context.Context, slider dtos.CreateSliderReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(slider); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	newSlider := models.Slider{
		ImagePathTM: slider.ImagePathTM,
		ImagePathEN: slider.ImagePathEN,
		ImagePathRU: slider.ImagePathRU,
		Platform:    slider.Platform,
	}

	brandID, err := s.repo.CreateSlider(ctx, newSlider)
	if err != nil {
		s.logger.Errorf("create err: %v", err)
		return brandID, err
	}
	return brandID, nil
}

func (s *SlidersService) GetAllSliders(ctx context.Context, limit, page int64, platform string) (dtos.SliderResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	sliders, count, err := s.repo.GetAllSliders(ctx, limit, offset, platform)
	if err != nil {
		s.logger.Errorf("get sliders err: %v", err)
		return dtos.SliderResult{}, err
	}
	var dtoSliders []dtos.Slider
	for _, b := range sliders {
		dtoSliders = append(dtoSliders, dtos.Slider{
			ID:          b.ID,
			ImagePathTM: b.ImagePathTM,
			ImagePathEN: b.ImagePathEN,
			ImagePathRU: b.ImagePathRU,
			Platform:    b.Platform,
		})
	}

	result := dtos.SliderResult{
		Sliders: dtoSliders,
		Count:   count,
	}
	return result, nil
}

func (s *SlidersService) UpdateSlider(ctx context.Context, slider dtos.UpdateSliderReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(slider); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	oldSlider, err := s.repo.GetSliderByID(ctx, slider.ID)
	if err != nil {
		s.logger.Errorf("get old slider err: %v", err)
		return 0, err
	}

	if oldSlider.ImagePathTM != slider.ImagePathTM && slider.ImagePathTM != "" {
		if err = helpers.DeleteImage(oldSlider.ImagePathTM); err != nil {
			s.logger.Errorf("delete old image path tm err: %v", err)
		}
	}

	if oldSlider.ImagePathEN != slider.ImagePathEN && slider.ImagePathEN != "" {
		if err = helpers.DeleteImage(oldSlider.ImagePathEN); err != nil {
			s.logger.Errorf("delete old image path en err: %v", err)
		}
	}

	if oldSlider.ImagePathRU != slider.ImagePathRU && slider.ImagePathRU != "" {
		if err = helpers.DeleteImage(oldSlider.ImagePathRU); err != nil {
			s.logger.Errorf("delete old image path ru err: %v", err)
		}
	}

	newSlider := models.Slider{
		ID:          slider.ID,
		ImagePathTM: slider.ImagePathTM,
		ImagePathEN: slider.ImagePathEN,
		ImagePathRU: slider.ImagePathRU,
		Platform:    slider.Platform,
	}

	sliderID, err := s.repo.UpdateSlider(ctx, newSlider)
	if err != nil {
		s.logger.Errorf("update slider err: %v", err)
		return sliderID, err
	}
	return sliderID, nil
}

func (s *SlidersService) DeleteSlider(ctx context.Context, id int64) error {
	oldSlider, err := s.repo.GetSliderByID(ctx, id)
	if err != nil {
		s.logger.Errorf("get old slider err: %v", err)
		return err
	}

	if oldSlider.ImagePathTM != "" {
		if err = helpers.DeleteImage(oldSlider.ImagePathTM); err != nil {
			s.logger.Errorf("delete old image path tm err: %v", err)
		}
	}

	if oldSlider.ImagePathEN != "" {
		if err = helpers.DeleteImage(oldSlider.ImagePathEN); err != nil {
			s.logger.Errorf("delete old image path en err: %v", err)
		}
	}

	if oldSlider.ImagePathRU != "" {
		if err = helpers.DeleteImage(oldSlider.ImagePathRU); err != nil {
			s.logger.Errorf("delete old image path ru err: %v", err)
		}
	}

	deleteID := models.ID{
		ID: id,
	}

	err = s.repo.DeleteSlider(ctx, deleteID)
	if err != nil {
		s.logger.Errorf("delete slider err: %v", err)
		return err
	}
	return nil
}
