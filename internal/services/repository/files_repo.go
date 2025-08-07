package repository

import (
	"autotm-admin/internal/dtos"
	"context"
	"mime/multipart"
)

type FilesService interface {
	UploadImage(file multipart.File, header *multipart.FileHeader) (string, error)
	DeleteImage(ctx context.Context, imagePath dtos.ImagePath) error
}
