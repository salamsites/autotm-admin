package services

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"autotm-admin/internal/repository/storage"
	"autotm-admin/internal/services/repository"
	"context"
	"fmt"
	"time"

	pushPb "autotm-admin/push_service_pb"

	"github.com/Hajymuhammet/elasticsearch-package/filter"
	"github.com/Hajymuhammet/elasticsearch-package/index"
	ms "github.com/Hajymuhammet/elasticsearch-package/models"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/elastic/go-elasticsearch/v8"
	slog "github.com/salamsites/package-log"
	spsql "github.com/salamsites/package-psql"
)

type CarsService struct {
	logger      *slog.Logger
	clientPsql  spsql.Client
	repo        storage.CarsRepository
	userService repository.UserService
	pushService pushPb.PushServiceClient
	stockRepo   storage.StockRepository
	esClient    *elasticsearch.Client
}

func NewCarsService(logger *slog.Logger, clientPsql spsql.Client, repo storage.CarsRepository, userService repository.UserService, pushService pushPb.PushServiceClient, stockRepo storage.StockRepository, esClient *elasticsearch.Client) *CarsService {
	return &CarsService{
		logger:      logger,
		clientPsql:  clientPsql,
		repo:        repo,
		userService: userService,
		pushService: pushService,
		stockRepo:   stockRepo,
		esClient:    esClient,
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
			DriveType:      car.DriveType,
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
			CreatedAt:      car.CreatedAt,
			UpdatedAt:      car.UpdatedAt,
		})
	}

	resp := dtos.CarsResp{
		Cars:  dtoCars,
		Count: count,
	}

	return resp, nil
}

func (s *CarsService) GetCarByID(ctx context.Context, id int64) (dtos.Car, error) {
	car, err := s.repo.GetCarByID(ctx, id)
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
		DriveType:      car.DriveType,
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
		CreatedAt:      car.CreatedAt,
		UpdatedAt:      car.UpdatedAt,
	}

	return result, nil
}

func (s *CarsService) UpdateCarStatus(ctx context.Context, req dtos.UpdateCarStatus) (dtos.ID, error) {
	var (
		id       dtos.ID
		carModel *ms.Car
	)
	trManager := manager.Must(trmpgx.NewDefaultFactory(s.clientPsql.Pool()))

	err := trManager.Do(ctx, func(ctx context.Context) error {

		validate := helpers.GetValidator()
		if err := validate.Struct(req); err != nil {
			s.logger.Errorf("validate err: %v", err)
			return err
		}

		carId, err := s.repo.UpdateCarStatus(ctx, req.ID, req.Status)
		if err != nil {
			s.logger.Errorf("update car status err: %v", err)
			return err
		}

		car, err := s.repo.GetCarByID(ctx, carId)
		if err != nil {
			s.logger.Errorf("get cars by id err: %w", err)
			return err
		}

		carModel = &ms.Car{
			ID:             car.Id,
			UserId:         car.UserId,
			UserName:       car.UserName,
			StockId:        car.StockId,
			StoreName:      car.StoreName,
			BrandId:        car.BrandId,
			BrandName:      car.BrandName,
			ModelId:        car.ModelId,
			ModelName:      car.ModelName,
			Year:           car.Year,
			Price:          car.Price,
			Color:          car.Color,
			Vin:            car.Vin,
			Description:    car.Description,
			CityId:         car.CityId,
			CityNameTM:     car.CityNameTM,
			CityNameEN:     car.CityNameEN,
			CityNameRU:     car.CityNameRU,
			Name:           car.Name,
			Mail:           car.Mail,
			PhoneNumber:    car.PhoneNumber,
			IsComment:      car.IsComment,
			IsExchange:     car.IsExchange,
			IsCredit:       car.IsCredit,
			Images:         car.Images,
			Status:         car.Status,
			Mileage:        car.Mileage,
			EngineCapacity: car.EngineCapacity,
			EngineType:     car.EngineType,
			BodyId:         car.BodyId,
			BodyNameTM:     car.BodyNameTM,
			BodyNameEN:     car.BodyNameEN,
			BodyNameRU:     car.BodyNameRU,
			Transmission:   car.Transmission,
			DriveType:      car.DriveType,
			CreatedAt:      car.CreatedAt,
			UpdatedAt:      car.UpdatedAt,
		}

		id.ID = carId

		return nil
	})

	if err != nil {
		s.logger.Errorf("update id car err: %v", err)
		return id, err
	}

	if err := index.UpdateCar(s.esClient, helpers.CarIndexName, carModel); err != nil {
		s.logger.Errorf("ES index error: %v", err)
	} else {
		s.logger.Infof("Car indexed successfully in Elasticsearch")
	}
	
	if req.StockID != 0 {
		if err := s.handlePushNotifications(req.StockID, req.Message); err != nil {
			s.logger.Errorf("push notification car error: %w", err)
		}
	}

	return id, nil
}

func (s *CarsService) GetTrucks(ctx context.Context, limit, page int64, search, status string) (dtos.TrucksResp, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	trucks, count, err := s.repo.GetTrucks(ctx, limit, offset, search, status)
	if err != nil {
		s.logger.Errorf("get trucks err: %v", err)
		return dtos.TrucksResp{}, err
	}

	var dtoTrucks []dtos.Truck
	for _, truck := range trucks {
		dtoTrucks = append(dtoTrucks, dtos.Truck{
			Id:              truck.Id,
			UserId:          truck.UserId,
			UserName:        truck.UserName,
			StockId:         truck.StockId,
			StoreName:       truck.StoreName,
			BrandId:         truck.BrandId,
			BrandName:       truck.BrandName,
			LoadCapacity:    truck.LoadCapacity,
			Price:           truck.Price,
			BodyType:        truck.BodyType,
			DriveType:       truck.DriveType,
			Transmission:    truck.Transmission,
			EngineType:      truck.EngineType,
			ModelId:         truck.ModelId,
			ModelName:       truck.ModelName,
			Year:            truck.Year,
			Seats:           truck.Seats,
			CabType:         truck.CabType,
			WheelFormula:    truck.WheelFormula,
			Chassis:         truck.Chassis,
			CabSuspension:   truck.CabSuspension,
			BusType:         truck.BusType,
			SuspensionType:  truck.SuspensionType,
			Brakes:          truck.Brakes,
			Axles:           truck.Axles,
			EngineHours:     truck.EngineHours,
			VehicleType:     truck.VehicleType,
			EngineCapacity:  truck.EngineCapacity,
			ForkliftType:    truck.ForkliftType,
			LiftingCapacity: truck.LiftingCapacity,
			Mileage:         truck.Mileage,
			ExcavatorType:   truck.ExcavatorType,
			BulldozerType:   truck.BulldozerType,
			Color:           truck.Color,
			Vin:             truck.Vin,
			BodyId:          truck.BodyId,
			BodyNameTM:      truck.BodyNameTM,
			BodyNameEN:      truck.BodyNameEN,
			BodyNameRU:      truck.BodyNameRU,
			Description:     truck.Description,
			CityId:          truck.CityId,
			CityNameTM:      truck.CityNameTM,
			CityNameEN:      truck.CityNameEN,
			CityNameRU:      truck.CityNameRU,
			Name:            truck.Name,
			Mail:            truck.Mail,
			PhoneNumber:     truck.PhoneNumber,
			IsComment:       truck.IsComment,
			IsExchange:      truck.IsExchange,
			IsCredit:        truck.IsCredit,
			Images:          truck.Images,
			Status:          truck.Status,
			CreatedAt:       truck.CreatedAt,
			UpdatedAt:       truck.UpdatedAt,
		})
	}

	resp := dtos.TrucksResp{
		Trucks: dtoTrucks,
		Count:  count,
	}

	return resp, nil
}

func (s *CarsService) GetTruckByID(ctx context.Context, id int64) (dtos.Truck, error) {
	truck, err := s.repo.GetTruckByID(ctx, id)
	if err != nil {
		s.logger.Errorf("get truck by id err: %v", err)
		return dtos.Truck{}, err
	}

	result := dtos.Truck{
		Id:              truck.Id,
		UserId:          truck.UserId,
		UserName:        truck.UserName,
		StockId:         truck.StockId,
		StoreName:       truck.StoreName,
		BrandId:         truck.BrandId,
		BrandName:       truck.BrandName,
		LoadCapacity:    truck.LoadCapacity,
		Price:           truck.Price,
		BodyType:        truck.BodyType,
		DriveType:       truck.DriveType,
		Transmission:    truck.Transmission,
		EngineType:      truck.EngineType,
		ModelId:         truck.ModelId,
		ModelName:       truck.ModelName,
		Year:            truck.Year,
		Seats:           truck.Seats,
		CabType:         truck.CabType,
		WheelFormula:    truck.WheelFormula,
		Chassis:         truck.Chassis,
		CabSuspension:   truck.CabSuspension,
		BusType:         truck.BusType,
		SuspensionType:  truck.SuspensionType,
		Brakes:          truck.Brakes,
		Axles:           truck.Axles,
		EngineHours:     truck.EngineHours,
		VehicleType:     truck.VehicleType,
		EngineCapacity:  truck.EngineCapacity,
		ForkliftType:    truck.ForkliftType,
		LiftingCapacity: truck.LiftingCapacity,
		Mileage:         truck.Mileage,
		ExcavatorType:   truck.ExcavatorType,
		BulldozerType:   truck.BulldozerType,
		Color:           truck.Color,
		Vin:             truck.Vin,
		BodyId:          truck.BodyId,
		BodyNameTM:      truck.BodyNameTM,
		BodyNameEN:      truck.BodyNameEN,
		BodyNameRU:      truck.BodyNameRU,
		Description:     truck.Description,
		CityId:          truck.CityId,
		CityNameTM:      truck.CityNameTM,
		CityNameEN:      truck.CityNameEN,
		CityNameRU:      truck.CityNameRU,
		Name:            truck.Name,
		Mail:            truck.Mail,
		PhoneNumber:     truck.PhoneNumber,
		IsComment:       truck.IsComment,
		IsExchange:      truck.IsExchange,
		IsCredit:        truck.IsCredit,
		Images:          truck.Images,
		Status:          truck.Status,
		CreatedAt:       truck.CreatedAt,
		UpdatedAt:       truck.UpdatedAt,
	}

	return result, nil
}

func (s *CarsService) UpdateTruckStatus(ctx context.Context, req dtos.UpdateTruckStatus) (dtos.ID, error) {
	var (
		id         dtos.ID
		truckModel *ms.Truck
	)

	trManager := manager.Must(trmpgx.NewDefaultFactory(s.clientPsql.Pool()))

	err := trManager.Do(ctx, func(ctx context.Context) error {
		validate := helpers.GetValidator()
		if err := validate.Struct(req); err != nil {
			s.logger.Errorf("validate err: %v", err)
			return err
		}

		truckId, err := s.repo.UpdateTruckStatus(ctx, req.ID, req.Status)
		if err != nil {
			s.logger.Errorf("update truck status err: %v", err)
			return err
		}

		truck, err := s.repo.GetTruckByID(ctx, truckId)
		if err != nil {
			s.logger.Errorf("get truck by id err: %w", err)
			return err
		}

		truckModel = &ms.Truck{
			Id:              truck.Id,
			UserId:          truck.UserId,
			UserName:        truck.UserName,
			StockId:         truck.StockId,
			StoreName:       truck.StoreName,
			BrandId:         truck.BrandId,
			BrandName:       truck.BrandName,
			LoadCapacity:    truck.LoadCapacity,
			Price:           truck.Price,
			BodyType:        truck.BodyType,
			DriveType:       truck.DriveType,
			Transmission:    truck.Transmission,
			EngineType:      truck.EngineType,
			ModelId:         truck.ModelId,
			ModelName:       truck.ModelName,
			Year:            truck.Year,
			Seats:           truck.Seats,
			CabType:         truck.CabType,
			WheelFormula:    truck.WheelFormula,
			Chassis:         truck.Chassis,
			CabSuspension:   truck.CabSuspension,
			BusType:         truck.BusType,
			SuspensionType:  truck.SuspensionType,
			Brakes:          truck.Brakes,
			Axles:           truck.Axles,
			EngineHours:     truck.EngineHours,
			VehicleType:     truck.VehicleType,
			EngineCapacity:  truck.EngineCapacity,
			ForkliftType:    truck.ForkliftType,
			LiftingCapacity: truck.LiftingCapacity,
			Mileage:         truck.Mileage,
			ExcavatorType:   truck.ExcavatorType,
			BulldozerType:   truck.BulldozerType,
			Color:           truck.Color,
			Vin:             truck.Vin,
			BodyId:          truck.BodyId,
			BodyNameTM:      truck.BodyNameTM,
			BodyNameEN:      truck.BodyNameEN,
			BodyNameRU:      truck.BodyNameRU,
			Description:     truck.Description,
			CityId:          truck.CityId,
			CityNameTM:      truck.CityNameTM,
			CityNameEN:      truck.CityNameEN,
			CityNameRU:      truck.CityNameRU,
			Name:            truck.Name,
			Mail:            truck.Mail,
			PhoneNumber:     truck.PhoneNumber,
			IsComment:       truck.IsComment,
			IsExchange:      truck.IsExchange,
			IsCredit:        truck.IsCredit,
			Images:          truck.Images,
			Status:          truck.Status,
			CreatedAt:       truck.CreatedAt,
			UpdatedAt:       truck.UpdatedAt,
		}

		id.ID = truckId

		return nil
	})

	if err != nil {
		s.logger.Errorf("update id truck err: %v", err)
		return id, err
	}

	if err := index.UpdateTruck(s.esClient, helpers.TruckIndexName, truckModel); err != nil {
		s.logger.Errorf("ES index error: %v", err)
	} else {
		s.logger.Infof("Truck indexed successfully in Elasticsearch")
	}

	if req.StockID != 0 {
		if err := s.handlePushNotifications(req.StockID, req.Message); err != nil {
			s.logger.Errorf("push notification truck error: %w", err)
		}
	}

	return id, nil
}

func (s *CarsService) GetMotors(ctx context.Context, limit, page int64, search, status string) (dtos.MotoResp, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	motors, count, err := s.repo.GetMotors(ctx, limit, offset, search, status)
	if err != nil {
		s.logger.Errorf("get motos err: %v", err)
		return dtos.MotoResp{}, err
	}

	var dtoMotors []dtos.Moto
	for _, moto := range motors {
		dtoMotors = append(dtoMotors, dtos.Moto{
			Id:                  moto.Id,
			UserId:              moto.UserId,
			UserName:            moto.UserName,
			StockId:             moto.StockId,
			StoreName:           moto.StoreName,
			BodyId:              moto.BodyId,
			BodyNameTM:          moto.BodyNameTM,
			BodyNameEN:          moto.BodyNameEN,
			BodyNameRU:          moto.BodyNameRU,
			BrandId:             moto.BrandId,
			BrandName:           moto.BrandName,
			ModelId:             moto.ModelId,
			ModelName:           moto.ModelName,
			TypeMotorcycles:     moto.TypeMotorcycles,
			Year:                moto.Year,
			Price:               moto.Price,
			Volume:              moto.Volume,
			EngineType:          moto.EngineType,
			NumberOfClockCycles: moto.NumberOfClockCycles,
			Mileage:             moto.Mileage,
			AirType:             moto.AirType,
			Color:               moto.Color,
			Vin:                 moto.Vin,
			Description:         moto.Description,
			CityId:              moto.CityId,
			CityNameTM:          moto.CityNameTM,
			CityNameEN:          moto.CityNameEN,
			CityNameRU:          moto.CityNameRU,
			Name:                moto.Name,
			Mail:                moto.Mail,
			PhoneNumber:         moto.PhoneNumber,
			IsComment:           moto.IsComment,
			IsExchange:          moto.IsExchange,
			IsCredit:            moto.IsCredit,
			Images:              moto.Images,
			Status:              moto.Status,
			Options:             moto.Options,
			CreatedAt:           moto.CreatedAt,
			UpdatedAt:           moto.UpdatedAt,
		})
	}

	resp := dtos.MotoResp{
		Motors: dtoMotors,
		Count:  count,
	}

	return resp, nil
}

func (s *CarsService) GetMotoByID(ctx context.Context, id int64) (dtos.Moto, error) {
	moto, err := s.repo.GetMotoByID(ctx, id)
	if err != nil {
		s.logger.Errorf("get moto by id err: %v", err)
		return dtos.Moto{}, err
	}

	result := dtos.Moto{
		Id:                  moto.Id,
		UserId:              moto.UserId,
		UserName:            moto.UserName,
		StockId:             moto.StockId,
		StoreName:           moto.StoreName,
		BodyId:              moto.BodyId,
		BodyNameTM:          moto.BodyNameTM,
		BodyNameEN:          moto.BodyNameEN,
		BodyNameRU:          moto.BodyNameRU,
		BrandId:             moto.BrandId,
		BrandName:           moto.BrandName,
		ModelId:             moto.ModelId,
		ModelName:           moto.ModelName,
		TypeMotorcycles:     moto.TypeMotorcycles,
		Year:                moto.Year,
		Price:               moto.Price,
		Volume:              moto.Volume,
		EngineType:          moto.EngineType,
		NumberOfClockCycles: moto.NumberOfClockCycles,
		Mileage:             moto.Mileage,
		AirType:             moto.AirType,
		Color:               moto.Color,
		Vin:                 moto.Vin,
		Description:         moto.Description,
		CityId:              moto.CityId,
		CityNameTM:          moto.CityNameTM,
		CityNameEN:          moto.CityNameEN,
		CityNameRU:          moto.CityNameRU,
		Name:                moto.Name,
		Mail:                moto.Mail,
		PhoneNumber:         moto.PhoneNumber,
		IsComment:           moto.IsComment,
		IsExchange:          moto.IsExchange,
		IsCredit:            moto.IsCredit,
		Images:              moto.Images,
		Status:              moto.Status,
		Options:             moto.Options,
		CreatedAt:           moto.CreatedAt,
		UpdatedAt:           moto.UpdatedAt,
	}

	return result, nil
}

func (s *CarsService) UpdateMotoStatus(ctx context.Context, req dtos.UpdateMotoStatus) (dtos.ID, error) {
	var (
		id        dtos.ID
		motoModel *ms.Moto
	)

	trManager := manager.Must(trmpgx.NewDefaultFactory(s.clientPsql.Pool()))

	err := trManager.Do(ctx, func(ctx context.Context) error {
		validate := helpers.GetValidator()
		if err := validate.Struct(req); err != nil {
			s.logger.Errorf("validation error: %v", err)
			return err
		}

		// Update status in DB
		motoId, err := s.repo.UpdateMotoStatus(ctx, req.ID, req.Status)
		if err != nil {
			s.logger.Errorf("update moto status error: %v", err)
			return err
		}

		// Get updated moto from DB
		moto, err := s.repo.GetMotoByID(ctx, motoId)
		if err != nil {
			s.logger.Errorf("get moto by id error: %v", err)
			return err
		}

		// Prepare data for Elasticsearch (outside the transaction)
		motoModel = &ms.Moto{
			Id:                  moto.Id,
			UserId:              moto.UserId,
			UserName:            moto.UserName,
			StockId:             moto.StockId,
			StoreName:           moto.StoreName,
			BodyId:              moto.BodyId,
			BodyNameTM:          moto.BodyNameTM,
			BodyNameEN:          moto.BodyNameEN,
			BodyNameRU:          moto.BodyNameRU,
			BrandId:             moto.BrandId,
			BrandName:           moto.BrandName,
			ModelId:             moto.ModelId,
			ModelName:           moto.ModelName,
			TypeMotorcycles:     moto.TypeMotorcycles,
			Year:                moto.Year,
			Price:               moto.Price,
			Volume:              moto.Volume,
			EngineType:          moto.EngineType,
			NumberOfClockCycles: moto.NumberOfClockCycles,
			Mileage:             moto.Mileage,
			AirType:             moto.AirType,
			Color:               moto.Color,
			Vin:                 moto.Vin,
			Description:         moto.Description,
			CityId:              moto.CityId,
			CityNameTM:          moto.CityNameTM,
			CityNameEN:          moto.CityNameEN,
			CityNameRU:          moto.CityNameRU,
			Name:                moto.Name,
			Mail:                moto.Mail,
			PhoneNumber:         moto.PhoneNumber,
			IsComment:           moto.IsComment,
			IsExchange:          moto.IsExchange,
			IsCredit:            moto.IsCredit,
			Images:              moto.Images,
			Status:              moto.Status,
			Options:             moto.Options,
			CreatedAt:           moto.CreatedAt,
			UpdatedAt:           moto.UpdatedAt,
		}

		id.ID = motoId
		return nil
	})

	if err != nil {
		s.logger.Errorf("update moto transaction failed: %v", err)
		return id, err
	}

	// Elasticsearch update
	if err := index.UpdateMoto(s.esClient, helpers.MotoIndexName, motoModel); err != nil {
		s.logger.Errorf("Elasticsearch update failed: %v", err)
	} else {
		s.logger.Infof("Moto indexed successfully in Elasticsearch")
	}

	// Push notification
	if req.StockID != 0 {
		if err := s.handlePushNotifications(req.StockID, req.Message); err != nil {
			s.logger.Errorf("push notification error: %v", err)
		}
	}

	return id, nil
}

func (s *CarsService) handlePushNotifications(stockID int64, message string) error {
	ctx := context.Background()
	const maxRetries = 3
	retryDelay := time.Second * 2

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := s.sendPushNotifications(ctx, stockID, message)
		if err == nil {
			s.logger.Infof("Push notifications sent successfully for stock %d", stockID)
			return err
		}

		s.logger.Warnf("Push attempt %d failed for stock %d: %v", attempt, stockID, err)

		if attempt == maxRetries {
			s.logger.Errorf("All push attempts failed for stock %d: %v", stockID, err)
			return err
		}

		time.Sleep(retryDelay)
		retryDelay *= 2 // Exponential backoff
	}
	return nil
}

func (s *CarsService) sendPushNotifications(ctx context.Context, stockID int64, message string) error {
	userIDs, err := s.stockRepo.GetStockFollowers(ctx, stockID)
	if err != nil {
		return fmt.Errorf("get stock followers: %w", err)
	}

	if len(userIDs) == 0 {
		s.logger.Debugf("No followers for stock %d", stockID)
		return nil
	}

	tokens, err := s.userService.GetUserFirebaseToken(ctx, userIDs)
	if err != nil {
		return fmt.Errorf("get firebase tokens: %w", err)
	}

	if len(tokens) == 0 {
		s.logger.Debugf("No tokens for stock %d followers", stockID)
		return nil
	}

	_, err = s.pushService.SendMultiPush(ctx, &pushPb.SendMultiPushRequest{
		Message: message,
		Tokens:  tokens,
	})
	if err != nil {
		return fmt.Errorf("send push: %w", err)
	}

	return nil
}

func (s *CarsService) SearchCars(ctx context.Context, ft *filter.CarFilter) ([]dtos.Car, error) {
	cars, err := filter.SearchCars(s.esClient, helpers.CarIndexName, ft)
	if err != nil {
		s.logger.Errorf("search cars ES err: %v", err)
		return nil, err
	}

	dtoCars := make([]dtos.Car, len(cars))

	for i, car := range cars {
		dtoCars[i] = dtos.Car{
			Id:             car.ID,
			UserId:         car.UserId,
			UserName:       car.UserName,
			StockId:        car.StockId,
			StoreName:      car.StoreName,
			BrandId:        car.BrandId,
			BrandName:      car.BrandName,
			ModelId:        car.ModelId,
			ModelName:      car.ModelName,
			Year:           car.Year,
			Price:          car.Price,
			Color:          car.Color,
			Vin:            car.Vin,
			Description:    car.Description,
			CityId:         car.CityId,
			CityNameTM:     car.CityNameTM,
			CityNameEN:     car.CityNameEN,
			CityNameRU:     car.CityNameRU,
			Name:           car.Name,
			Mail:           car.Mail,
			PhoneNumber:    car.PhoneNumber,
			IsComment:      car.IsComment,
			IsExchange:     car.IsExchange,
			IsCredit:       car.IsCredit,
			Images:         car.Images,
			Status:         car.Status,
			Mileage:        car.Mileage,
			EngineCapacity: car.EngineCapacity,
			EngineType:     car.EngineType,
			BodyId:         car.BodyId,
			BodyNameTM:     car.BodyNameTM,
			BodyNameEN:     car.BodyNameEN,
			BodyNameRU:     car.BodyNameRU,
			Transmission:   car.Transmission,
			DriveType:      car.DriveType,
			CreatedAt:      car.CreatedAt,
			UpdatedAt:      car.UpdatedAt,
		}
	}

	return dtoCars, nil
}
