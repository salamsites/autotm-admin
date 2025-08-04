package services

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"context"
	slog "github.com/salamsites/package-log"
	"mime/multipart"
)

type FilesServiceImpl struct {
	logger *slog.Logger
}

func NewFilesService(logger *slog.Logger) *FilesServiceImpl {
	return &FilesServiceImpl{
		logger: logger,
	}
}

func (s *FilesServiceImpl) UploadImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	imagePath, err := helpers.UploadImage(file, header)
	if err != nil {
		s.logger.Errorf("upload image err: %v", err)
		return imagePath, err
	}
	return imagePath, nil
}

func (s *FilesServiceImpl) DeleteImage(ctx context.Context, imagePath dtos.ImagePath) error {
	if imagePath.ImagePath != "" {
		if err := helpers.DeleteImage(imagePath.ImagePath); err != nil {
			s.logger.Errorf("delete image err: %v", err)
		}
	}

	if len(imagePath.Images) > 0 {
		for _, img := range imagePath.Images {
			if img == "" {
				continue
			}
			if err := helpers.DeleteImage(img); err != nil {
				s.logger.Errorf("delete image err: %v", err)
				return err
			}
		}
	}

	return nil
}
