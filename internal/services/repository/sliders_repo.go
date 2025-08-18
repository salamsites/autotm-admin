package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type SlidersService interface {
	CreateSlider(ctx context.Context, slider dtos.CreateSliderReq) (dtos.ID, error)
	GetAllSliders(ctx context.Context, limit, page int64, platform string) (dtos.SliderResult, error)
	UpdateSlider(ctx context.Context, role dtos.UpdateSliderReq) (dtos.ID, error)
	DeleteSlider(ctx context.Context, id int64) error
}
