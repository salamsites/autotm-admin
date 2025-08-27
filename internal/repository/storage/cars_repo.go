package storage

import (
	"autotm-admin/internal/models"
	"context"
)

type CarsRepository interface {
	GetCars(ctx context.Context, limit, page int64, search, status string) ([]models.Car, int64, error)
	GetCarsByID(ctx context.Context, id int64) (models.Car, error)
	UpdateCarStatus(ctx context.Context, id int64, status string) (int64, error)
	GetUserByCarId(ctx context.Context, carId int64) (int64, error)

	//Truck
	GetTrucks(ctx context.Context, limit, page int64, search, status string) ([]models.Truck, int64, error)
}
