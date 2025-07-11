package storage

import (
	"autotm-admin/internal/models"
	"context"
)

type SlidersRepository interface {
	CreateSlider(ctx context.Context, slider models.Slider) (int64, error)
	GetAllSliders(ctx context.Context, limit, page int64, search string) ([]models.Slider, int64, error)
	UpdateSlider(ctx context.Context, slider models.Slider) (int64, error)
	GetSliderByID(ctx context.Context, id int64) (models.Slider, error)
	DeleteSlider(ctx context.Context, id models.ID) error
}
