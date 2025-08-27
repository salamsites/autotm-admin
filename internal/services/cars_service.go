package services

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"autotm-admin/internal/repository/storage"
	"autotm-admin/internal/services/repository"
	"context"

	slog "github.com/salamsites/package-log"
)

type CarsService struct {
	logger      *slog.Logger
	repo        storage.CarsRepository
	userService repository.UserService
	pushService repository.PushService
}

func NewCarsService(logger *slog.Logger, repo storage.CarsRepository, userService repository.UserService, pushService repository.PushService) *CarsService {
	return &CarsService{
		logger:      logger,
		repo:        repo,
		userService: userService,
		pushService: pushService,
	}
}

func (s *CarsService) GetCars(ctx context.Context, limit, page int64, search, status string) (dtos.CarsResp, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	cars, count, err := s.repo.GetCars(ctx, limit, offset, search, status)
	if err != nil {
		s.logger.Errorf("get cars err: %v", err)
		return dtos.CarsResp{}, err
	}

	var dtoCars []dtos.Car
	for _, car := range cars {
		dtoCars = append(dtoCars, dtos.Car{
			Id:             car.Id,
			UserId:         car.UserId,
			UserName:       car.UserName,
			StockId:        car.StockId,
			StoreName:      car.StoreName,
			BrandId:        car.BrandId,
			BrandName:      car.BrandName,
			ModelId:        car.ModelId,
			ModelName:      car.ModelName,
			Year:           car.Year,
			Mileage:        car.Mileage,
			Color:          car.Color,
			EngineCapacity: car.EngineCapacity,
			EngineType:     car.EngineType,
			BodyId:         car.BodyId,
			BodyNameTM:     car.BodyNameTM,
			BodyNameEN:     car.BodyNameEN,
			BodyNameRU:     car.BodyNameRU,
			Transmission:   car.Transmission,
			DriveType:      car.DriverType,
			Vin:            car.Vin,
			Description:    car.Description,
			CityId:         car.CityId,
			CityNameTM:     car.CityNameTM,
			CityNameEN:     car.CityNameEN,
			CityNameRU:     car.CityNameRU,
			Name:           car.Name,
			Mail:           car.Mail,
			PhoneNumber:    car.PhoneNumber,
			Price:          car.Price,
			IsComment:      car.IsComment,
			IsExchange:     car.IsExchange,
			IsCredit:       car.IsCredit,
			Images:         car.Images,
			Status:         car.Status,
		})
	}

	resp := dtos.CarsResp{
		Cars:  dtoCars,
		Count: count,
	}

	return resp, nil
}

func (s *CarsService) GetCarsByID(ctx context.Context, id int64) (dtos.Car, error) {
	car, err := s.repo.GetCarsByID(ctx, id)
	if err != nil {
		s.logger.Errorf("get cars by id err: %v", err)
		return dtos.Car{}, err
	}

	result := dtos.Car{
		Id:             car.Id,
		UserId:         car.UserId,
		UserName:       car.UserName,
		StockId:        car.StockId,
		StoreName:      car.StoreName,
		BrandId:        car.BrandId,
		BrandName:      car.BrandName,
		ModelId:        car.ModelId,
		ModelName:      car.ModelName,
		Year:           car.Year,
		Mileage:        car.Mileage,
		Color:          car.Color,
		EngineCapacity: car.EngineCapacity,
		EngineType:     car.EngineType,
		BodyId:         car.BodyId,
		BodyNameTM:     car.BodyNameTM,
		BodyNameEN:     car.BodyNameEN,
		BodyNameRU:     car.BodyNameRU,
		Transmission:   car.Transmission,
		DriveType:      car.DriverType,
		Vin:            car.Vin,
		Description:    car.Description,
		CityId:         car.CityId,
		CityNameTM:     car.CityNameTM,
		CityNameEN:     car.CityNameEN,
		CityNameRU:     car.CityNameRU,
		Name:           car.Name,
		Mail:           car.Mail,
		PhoneNumber:    car.PhoneNumber,
		Price:          car.Price,
		IsComment:      car.IsComment,
		IsExchange:     car.IsExchange,
		IsCredit:       car.IsCredit,
		Images:         car.Images,
		Status:         car.Status,
	}

	return result, nil
}

func (s *CarsService) UpdateCarStatus(ctx context.Context, car dtos.UpdateCarStatus) (dtos.ID, error) {
	var id dtos.ID

	validate := helpers.GetValidator()
	if err := validate.Struct(car); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	carId, err := s.repo.UpdateCarStatus(ctx, car.ID, car.Status)
	if err != nil {
		s.logger.Errorf("update car status err: %v", err)
		return id, err
	}

	userId, err := s.repo.GetUserByCarId(ctx, carId)
	if err != nil {
		s.logger.Errorf("get GetUserByCarId err: %v", err)
		return id, err
	}

	token, err := s.userService.GetUserFirebaseToken(ctx, userId)
	if err != nil {
		s.logger.Errorf("get GetUserFirebaseToken err: %v", err)
		return id, err
	}

	reqPush := dtos.ReqSendPushDTO{
		Message: car.Message,
		Token:   token,
	}

	go s.pushService.SendPush(reqPush)

	id.ID = carId
	return id, nil
}
