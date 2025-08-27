package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type CarsService interface {
	GetCars(ctx context.Context, limit, page int64, search, status string) (dtos.CarsResp, error)
	GetCarsByID(ctx context.Context, id int64) (dtos.Car, error)
	UpdateCarStatus(ctx context.Context, car dtos.UpdateCarStatus) (dtos.ID, error)
}
