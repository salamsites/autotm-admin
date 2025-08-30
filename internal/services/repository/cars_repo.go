package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type CarsService interface {
	//Cars
	GetCars(ctx context.Context, limit, page int64, search, status string) (dtos.CarsResp, error)
	GetCarByID(ctx context.Context, id int64) (dtos.Car, error)
	UpdateCarStatus(ctx context.Context, car dtos.UpdateCarStatus) (dtos.ID, error)

	//Trucks
	GetTrucks(ctx context.Context, limit, page int64, search, status string) (dtos.TrucksResp, error)
	GetTruckByID(ctx context.Context, id int64) (dtos.Truck, error)
	UpdateTruckStatus(ctx context.Context, truck dtos.UpdateTruckStatus) (dtos.ID, error)

	//Motors
	GetMotors(ctx context.Context, limit, page int64, search, status string) (dtos.MotoResp, error)
	GetMotoByID(ctx context.Context, id int64) (dtos.Moto, error)
	UpdateMotoStatus(ctx context.Context, moto dtos.UpdateMotoStatus) (dtos.ID, error)
}
