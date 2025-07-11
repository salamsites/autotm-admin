package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type SlidersService interface {
	CreateSlider(ctx context.Context, slider dtos.Slider) (int64, error)
	GetAllSliders(ctx context.Context, limit, page int64, platform string) (dtos.SliderResult, error)
	UpdateSlider(ctx context.Context, role dtos.Slider) (int64, error)
	DeleteSlider(ctx context.Context, id int64) error
}
