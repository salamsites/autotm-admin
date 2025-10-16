package services

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/models"
	"autotm-admin/internal/repository/storage"
	"context"

	sminio "github.com/salamsites/minio-pkg"
	slog "github.com/salamsites/package-log"
)

type SlidersService struct {
	logger          *slog.Logger
	repo            storage.SlidersRepository
	minioFileClient sminio.ImageClient
}

func NewSlidersService(logger *slog.Logger, repo storage.SlidersRepository, minioFileClient sminio.ImageClient) *SlidersService {
	return &SlidersService{
		logger:          logger,
		repo:            repo,
		minioFileClient: minioFileClient,
	}
}

func (s *SlidersService) CreateSlider(ctx context.Context, slider dtos.CreateSliderReq) (dtos.ID, error) {
	var id dtos.ID

	newSlider := models.Slider{
		ImagePathTM: slider.ImagePathTM,
		ImagePathEN: slider.ImagePathEN,
		ImagePathRU: slider.ImagePathRU,
		UploadIdTM:  slider.UploadIdTM,
		UploadIdEN:  slider.UploadIdEN,
		UploadIdRU:  slider.UploadIdRU,
	}

	brandID, err := s.repo.CreateSlider(ctx, newSlider)
	if err != nil {
		s.logger.Errorf("create slider service err: %v", err)
		return id, err
	}
	id.ID = brandID
	return id, nil
}

func (s *SlidersService) GetAllSliders(ctx context.Context, limit, page int64) (dtos.SliderResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	sliders, err := s.repo.GetAllSliders(ctx, limit, offset)
	if err != nil {
		s.logger.Errorf("get sliders service err: %v", err)
		return dtos.SliderResult{}, err
	}
	var dtoSliders []dtos.Slider
	for _, b := range sliders {
		dtoSliders = append(dtoSliders, dtos.Slider{
			ID:          b.ID,
			ImagePathTM: b.ImagePathTM,
			ImagePathEN: b.ImagePathEN,
			ImagePathRU: b.ImagePathRU,
			UploadIdTM:  b.UploadIdTM,
			UploadIdEN:  b.UploadIdEN,
			UploadIdRU:  b.UploadIdRU,
		})
	}

	count, errCount := s.repo.GetSlidersCount(ctx)
	if errCount != nil {
		s.logger.Errorf("get slidersCount service err: %v", err)
		return dtos.SliderResult{}, errCount
	}

	result := dtos.SliderResult{
		Sliders: dtoSliders,
		Count:   count,
	}
	return result, nil
}

func (s *SlidersService) UpdateSlider(ctx context.Context, slider dtos.UpdateSliderReq) (dtos.ID, error) {
	var id dtos.ID

	newSlider := models.Slider{
		ID:          slider.ID,
		ImagePathTM: slider.ImagePathTM,
		ImagePathEN: slider.ImagePathEN,
		ImagePathRU: slider.ImagePathRU,
		UploadIdTM:  slider.UploadIdTM,
		UploadIdEN:  slider.UploadIdEN,
		UploadIdRU:  slider.UploadIdRU,
	}

	sliderID, err := s.repo.UpdateSlider(ctx, newSlider)
	if err != nil {
		s.logger.Errorf("update slider service err: %v", err)
		return id, err
	}
	id.ID = sliderID
	return id, nil
}

func (s *SlidersService) DeleteSlider(ctx context.Context, id int64) error {
	err := s.repo.DeleteSlider(ctx, id)
	if err != nil {
		s.logger.Errorf("delete slider service err: %v", err)
		return err
	}
	return nil
}
