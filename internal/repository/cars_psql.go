package repository

import (
	"autotm-admin/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	slog "github.com/salamsites/package-log"
	spsql "github.com/salamsites/package-psql"
)

type CarsPsqlRepository struct {
	logger *slog.Logger
	client spsql.Client
}

func NewCarsPsqlRepository(logger *slog.Logger, client spsql.Client) *CarsPsqlRepository {
	return &CarsPsqlRepository{
		logger: logger,
		client: client,
	}
}

func (r *CarsPsqlRepository) GetCars(ctx context.Context, limit, page int64, search, status string) ([]models.Car, int64, error) {
	var (
		cars  []models.Car
		count int64
	)

	query := `
		SELECT
			cr.id, cr.user_id, u.full_name, cr.stock_id, s.store_name, cr.brand_id, b.name,
			cr.model_id, m.name, cr.year, cr.mileage, cr.color, cr.engine_capacity, cr.engine_type,
			cr.body_id, bt.name_tm, bt.name_en, bt.name_ru, cr.transmission, cr.drive_type, cr.vin, 
			cr.description, cr.city_id, cs.name_tm, cs.name_en, cs.name_ru, cr.name, cr.mail, cr.phone_number, 
			cr.price, cr.is_comment, cr.is_exchange, cr.is_credit, cr.images, cr.status
		FROM cars cr
			LEFT JOIN users u ON u.id = cr.user_id 
			LEFT JOIN stocks s ON s.id = cr.stock_id
			LEFT JOIN brands b ON b.id = cr.brand_id
			LEFT JOIN models m ON m.id = cr.model_id
			LEFT JOIN body_types bt ON bt.id = cr.body_id
			LEFT JOIN cities cs ON cs.id = cr.city_id
		WHERE (u.full_name ILIKE '%' || @search || '%' OR s.store_name ILIKE '%' || @search || '%')
	`

	args := pgx.NamedArgs{
		"search": search,
		"limit":  limit,
		"offset": page,
	}

	if status != "" {
		query += " AND cr.status = @status "
		args["status"] = status
	}

	query += `
   		ORDER BY cr.created_at DESC
   		LIMIT @limit OFFSET @offset
    `

	rows, err := r.client.Query(ctx, query, args)
	if err != nil {
		r.logger.Errorf("Error getting cars: %s", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var car models.Car

		err := rows.Scan(
			&car.Id,
			&car.UserId,
			&car.UserName,
			&car.StockId,
			&car.StoreName,
			&car.BrandId,
			&car.BrandName,
			&car.ModelId,
			&car.ModelName,
			&car.Year,
			&car.Mileage,
			&car.Color,
			&car.EngineCapacity,
			&car.EngineType,
			&car.BodyId,
			&car.BodyNameTM,
			&car.BodyNameEN,
			&car.BodyNameRU,
			&car.Transmission,
			&car.DriveType,
			&car.Vin,
			&car.Description,
			&car.CityId,
			&car.CityNameTM,
			&car.CityNameEN,
			&car.CityNameRU,
			&car.Name,
			&car.Mail,
			&car.PhoneNumber,
			&car.Price,
			&car.IsComment,
			&car.IsExchange,
			&car.IsCredit,
			&car.Images,
			&car.Status,
		)
		if err != nil {
			r.logger.Errorf("Error getting cars: %s", err)
		}
		cars = append(cars, car)
	}

	queryCount := `
			SELECT 
    			COUNT(cr.id) 
			FROM cars cr
				LEFT JOIN users u ON u.id = cr.user_id 
				LEFT JOIN stocks s ON s.id = cr.stock_id
				LEFT JOIN brands b ON b.id = cr.brand_id
				LEFT JOIN models m ON m.id = cr.model_id
				LEFT JOIN body_types bt ON bt.id = cr.body_id
				LEFT JOIN cities cs ON cs.id = cr.city_id
			WHERE (u.full_name ILIKE '%' || @search || '%' OR s.store_name ILIKE '%' || @search || '%')
		`
	argsCount := pgx.NamedArgs{
		"search": search,
	}

	if status != "" {
		queryCount += " AND cr.status = @status "
		argsCount["status"] = status
	}

	err = r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
	if err != nil {
		r.logger.Errorf("Error getting cars count: %s", err)
		return nil, 0, err
	}

	return cars, count, nil
}

func (r *CarsPsqlRepository) GetCarByID(ctx context.Context, id int64) (models.Car, error) {
	var car models.Car

	query := `
		SELECT
			cr.id, cr.user_id, u.full_name, cr.stock_id, s.store_name, cr.brand_id, b.name,
			cr.model_id, m.name, cr.year, cr.mileage, cr.color, cr.engine_capacity, cr.engine_type,
			cr.body_id, bt.name_tm, bt.name_en, bt.name_ru, cr.transmission, cr.drive_type, cr.vin, 
			cr.description, cr.city_id, cs.name_tm, cs.name_en, cs.name_ru, cr.name, cr.mail, cr.phone_number, 
			cr.price, cr.is_comment, cr.is_exchange, cr.is_credit, cr.images, cr.status
		FROM cars cr
			LEFT JOIN users u ON u.id = cr.user_id 
			LEFT JOIN stocks s ON s.id = cr.stock_id
			LEFT JOIN brands b ON b.id = cr.brand_id
			LEFT JOIN models m ON m.id = cr.model_id
			LEFT JOIN body_types bt ON bt.id = cr.body_id
			LEFT JOIN cities cs ON cs.id = cr.city_id
		WHERE cr.id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&car.Id, &car.UserId, &car.UserName, &car.StockId, &car.StoreName,
		&car.BrandId, &car.Name, &car.ModelId, &car.Name, &car.Name, &car.Year, &car.Mileage, &car.Color,
		&car.EngineCapacity, &car.EngineType, &car.BodyId, &car.BodyNameTM, &car.BodyNameEN, &car.BodyNameRU, &car.Transmission,
		&car.DriveType, &car.Vin, &car.Description, &car.CityId, &car.CityNameTM, &car.CityNameEN, &car.CityNameRU, &car.Name, &car.Mail,
		&car.PhoneNumber, &car.Price, &car.IsComment, &car.IsExchange, &car.IsCredit, &car.Images, &car.Status,
	)

	if err != nil {
		r.logger.Errorf("Error getting car by id: %s", err)
		return car, err
	}

	return car, nil
}

func (r *CarsPsqlRepository) UpdateCarStatus(ctx context.Context, id int64, status string) (int64, error) {
	var carId int64

	query := `
   		UPDATE cars SET
   		     status = @status
   		WHERE id = @id
   		RETURNING id
	`

	args := pgx.NamedArgs{
		"status": status,
		"id":     id,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&carId)
	if err != nil {
		r.logger.Errorf("update car status err: %v", err)
		return carId, err
	}

	return carId, nil
}

func (r *CarsPsqlRepository) GetUserByCarId(ctx context.Context, carId int64) (int64, error) {
	var userId int64

	query := `
		SELECT 
			user_id
		FROM cars
		WHERE id = @car_id
	`

	args := pgx.NamedArgs{
		"car_id": carId,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&userId)
	if err != nil {
		r.logger.Errorf("get user by car id err: %v", err)
		return userId, err
	}
	return userId, nil
}

func (r *CarsPsqlRepository) GetTrucks(ctx context.Context, limit, page int64, search, status string) ([]models.Truck, int64, error) {
	var (
		trucks []models.Truck
		count  int64
	)

	query := `
		SELECT
			t.id, t.user_id, u.full_name, t.stock_id, s.store_name, t.brand_id, b.name,
			t.load_capacity, t.price, t.body_type, t.drive_type, t.transmission, t.engine_type,
			t.model_id, m.name, t.year, t.seats, t.cab_type, t.wheel_formula, t.chassis, t.cab_suspension,
			t.bus_type, t.suspension_type, t.brakes, t.axles, t.engine_hours, t.vehicle_type, t.engine_capacity,
			t.forklift_type, t.lifting_capacity, t.mileage, t.excavator_type, t.bulldozer_type, t.color, t.vin, 
			t.body_id, bt.name_tm, bt.name_en, bt.name_ru, t.description, t.city_id, cs.name_tm, cs.name_en, cs.name_ru, 
			t.name, t.mail, t.phone_number, t.is_comment, t.is_exchange, t.is_credit, t.images, t.status
		FROM trucks t
			LEFT JOIN users u ON u.id = t.user_id 
			LEFT JOIN stocks s ON s.id = t.stock_id
			LEFT JOIN brands b ON b.id = t.brand_id
			LEFT JOIN models m ON m.id = t.model_id
			LEFT JOIN body_types bt ON bt.id = t.body_id
			LEFT JOIN cities cs ON cs.id = t.city_id
		WHERE (u.full_name ILIKE '%' || @search || '%' OR s.store_name ILIKE '%' || @search || '%')
	`

	args := pgx.NamedArgs{
		"search": search,
		"limit":  limit,
		"offset": page,
	}

	if status != "" {
		query += " AND t.status = @status "
		args["status"] = status
	}

	query += `
   		ORDER BY t.created_at DESC
   		LIMIT @limit OFFSET @offset
    `

	rows, err := r.client.Query(ctx, query, args)
	if err != nil {
		r.logger.Errorf("Error getting trucks: %s", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var truck models.Truck

		err := rows.Scan(
			&truck.Id,
			&truck.UserId,
			&truck.UserName,
			&truck.StockId,
			&truck.StoreName,
			&truck.BrandId,
			&truck.BrandName,
			&truck.LoadCapacity,
			&truck.Price,
			&truck.BodyType,
			&truck.DriveType,
			&truck.Transmission,
			&truck.EngineType,
			&truck.ModelId,
			&truck.ModelName,
			&truck.Year,
			&truck.Seats,
			&truck.CabType,
			&truck.WheelFormula,
			&truck.Chassis,
			&truck.CabSuspension,
			&truck.BusType,
			&truck.SuspensionType,
			&truck.Brakes,
			&truck.Axles,
			&truck.EngineHours,
			&truck.VehicleType,
			&truck.EngineCapacity,
			&truck.ForkliftType,
			&truck.LiftingCapacity,
			&truck.Mileage,
			&truck.ExcavatorType,
			&truck.BulldozerType,
			&truck.Color,
			&truck.Vin,
			&truck.BodyId,
			&truck.BodyNameTM,
			&truck.BodyNameEN,
			&truck.BodyNameRU,
			&truck.Description,
			&truck.CityId,
			&truck.CityNameTM,
			&truck.CityNameEN,
			&truck.CityNameRU,
			&truck.Name,
			&truck.Mail,
			&truck.PhoneNumber,
			&truck.IsComment,
			&truck.IsExchange,
			&truck.IsCredit,
			&truck.Images,
			&truck.Status,
		)
		if err != nil {
			r.logger.Errorf("Error getting cars: %s", err)
		}
		trucks = append(trucks, truck)
	}

	queryCount := `
			SELECT 
    			COUNT(t.id) 
			FROM trucks t
				LEFT JOIN users u ON u.id = t.user_id 
				LEFT JOIN stocks s ON s.id = t.stock_id
				LEFT JOIN brands b ON b.id = t.brand_id
				LEFT JOIN models m ON m.id = t.model_id
				LEFT JOIN body_types bt ON bt.id = t.body_id
				LEFT JOIN cities cs ON cs.id = t.city_id
			WHERE (u.full_name ILIKE '%' || @search || '%' OR s.store_name ILIKE '%' || @search || '%')
		`
	argsCount := pgx.NamedArgs{
		"search": search,
	}

	if status != "" {
		queryCount += " AND t.status = @status "
		argsCount["status"] = status
	}

	err = r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
	if err != nil {
		r.logger.Errorf("Error getting trucks count: %s", err)
		return nil, 0, err
	}

	return trucks, count, nil
}

func (r *CarsPsqlRepository) GetTruckByID(ctx context.Context, id int64) (models.Truck, error) {
	var truck models.Truck

	query := `
		SELECT
			t.id, t.user_id, u.full_name, t.stock_id, s.store_name, t.brand_id, b.name,
			t.load_capacity, t.price, t.body_type, t.drive_type, t.transmission, t.engine_type,
			t.model_id, m.name, t.year, t.seats, t.cab_type, t.wheel_formula, t.chassis, t.cab_suspension,
			t.bus_type, t.suspension_type, t.brakes, t.axles, t.engine_hours, t.vehicle_type, t.engine_capacity,
			t.forklift_type, t.lifting_capacity, t.mileage, t.excavator_type, t.bulldozer_type, t.color, t.vin, 
			t.body_id, bt.name_tm, bt.name_en, bt.name_ru, t.description, t.city_id, cs.name_tm, cs.name_en, cs.name_ru, 
			t.name, t.mail, t.phone_number, t.is_comment, t.is_exchange, t.is_credit, t.images, t.status
		FROM trucks t
			LEFT JOIN users u ON u.id = t.user_id 
			LEFT JOIN stocks s ON s.id = t.stock_id
			LEFT JOIN brands b ON b.id = t.brand_id
			LEFT JOIN models m ON m.id = t.model_id
			LEFT JOIN body_types bt ON bt.id = t.body_id
			LEFT JOIN cities cs ON cs.id = t.city_id
		WHERE t.id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&truck.Id, &truck.UserId, &truck.UserName, &truck.StockId, &truck.StoreName,
		&truck.BrandId, &truck.BrandName, &truck.LoadCapacity, &truck.Price, &truck.BodyType, &truck.DriveType, &truck.Transmission,
		&truck.EngineType, &truck.ModelId, &truck.ModelName, &truck.Year, &truck.Seats, &truck.CabType, &truck.WheelFormula, &truck.Chassis,
		&truck.CabSuspension, &truck.BusType, &truck.SuspensionType, &truck.Brakes, &truck.Axles, &truck.EngineHours, &truck.VehicleType,
		&truck.EngineCapacity, &truck.ForkliftType, &truck.LiftingCapacity, &truck.Mileage, &truck.ExcavatorType, &truck.BulldozerType,
		&truck.Color, &truck.Vin, &truck.BodyId, &truck.BodyNameTM, &truck.BodyNameEN, &truck.BodyNameRU, &truck.Description,
		&truck.CityId, &truck.CityNameTM, &truck.CityNameEN, &truck.CityNameRU, &truck.Name, &truck.Mail, &truck.PhoneNumber,
		&truck.IsComment, &truck.IsExchange, &truck.IsCredit, &truck.Images, &truck.Status,
	)

	if err != nil {
		r.logger.Errorf("Error getting truck by id: %s", err)
		return truck, err
	}

	return truck, nil
}

func (r *CarsPsqlRepository) UpdateTruckStatus(ctx context.Context, id int64, status string) (int64, error) {
	var truckId int64

	query := `
   		UPDATE trucks SET
   		     status = @status
   		WHERE id = @id
   		RETURNING id
	`

	args := pgx.NamedArgs{
		"status": status,
		"id":     id,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&truckId)
	if err != nil {
		r.logger.Errorf("update truck status err: %v", err)
		return truckId, err
	}

	return truckId, nil
}

func (r *CarsPsqlRepository) GetUserByTruckId(ctx context.Context, truckId int64) (int64, error) {
	var userId int64

	query := `
		SELECT 
			user_id
		FROM trucks
		WHERE id = @truck_id
	`

	args := pgx.NamedArgs{
		"truck_id": truckId,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&userId)
	if err != nil {
		r.logger.Errorf("get user by truck id err: %v", err)
		return userId, err
	}
	return userId, nil
}

func (r *CarsPsqlRepository) GetMotors(ctx context.Context, limit, page int64, search, status string) ([]models.Moto, int64, error) {
	var (
		motors []models.Moto
		count  int64
	)

	query := `
		SELECT
			ms.id, ms.user_id, u.full_name, ms.stock_id, s.store_name, ms.brand_id, b.name,
			ms.type_motorcycles, ms.year, ms.price, ms.volume, ms.engine_type, 
			ms.number_of_clock_cycles, t.model_id, m.name, ms.air_type, ms.color, ms.vin, 
			ms.description, ms.city_id, cs.name_tm, cs.name_en, cs.name_ru,
			ms.name, ms.mail, ms.phone_number, ms.options, ms.is_comment, 
			ms.is_exchange, ms.is_credit, ms.images, ms.status
		FROM motos ms
			LEFT JOIN users u ON u.id = ms.user_id 
			LEFT JOIN stocks s ON s.id = ms.stock_id
			LEFT JOIN brands b ON b.id = ms.brand_id
			LEFT JOIN models m ON m.id = ms.model_id
			LEFT JOIN body_types bt ON bt.id = ms.body_id
			LEFT JOIN cities cs ON cs.id = ms.city_id
		WHERE (u.full_name ILIKE '%' || @search || '%' OR s.store_name ILIKE '%' || @search || '%')
	`

	args := pgx.NamedArgs{
		"search": search,
		"limit":  limit,
		"offset": page,
	}

	if status != "" {
		query += " AND ms.status = @status "
		args["status"] = status
	}

	query += `
   		ORDER BY ms.created_at DESC
   		LIMIT @limit OFFSET @offset
    `

	rows, err := r.client.Query(ctx, query, args)
	if err != nil {
		r.logger.Errorf("Error getting motos: %s", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var motor models.Moto

		err := rows.Scan(
			&motor.Id,
			&motor.UserId,
			&motor.UserName,
			&motor.StockId,
			&motor.StoreName,
			&motor.BrandId,
			&motor.BrandName,
			&motor.TypeMotorcycles,
			&motor.Year,
			&motor.Price,
			&motor.Volume,
			&motor.EngineType,
			&motor.NumberOfClockCycles,
			&motor.ModelId,
			&motor.ModelName,
			&motor.AirType,
			&motor.Color,
			&motor.Vin,
			&motor.Description,
			&motor.CityId,
			&motor.CityNameTM,
			&motor.CityNameEN,
			&motor.CityNameRU,
			&motor.Name,
			&motor.Mail,
			&motor.PhoneNumber,
			&motor.Options,
			&motor.IsComment,
			&motor.IsExchange,
			&motor.IsCredit,
			&motor.Images,
			&motor.Status,
		)
		if err != nil {
			r.logger.Errorf("Error getting motors: %s", err)
		}
		motors = append(motors, motor)
	}

	queryCount := `
			SELECT 
    			COUNT(ms.id) 
			FROM motos ms
				LEFT JOIN users u ON u.id = ms.user_id 
				LEFT JOIN stocks s ON s.id = ms.stock_id
				LEFT JOIN brands b ON b.id = ms.brand_id
				LEFT JOIN models m ON m.id = ms.model_id
				LEFT JOIN body_types bt ON bt.id = ms.body_id
				LEFT JOIN cities cs ON cs.id = ms.city_id
			WHERE (u.full_name ILIKE '%' || @search || '%' OR s.store_name ILIKE '%' || @search || '%')
		`
	argsCount := pgx.NamedArgs{
		"search": search,
	}

	if status != "" {
		queryCount += " AND ms.status = @status "
		argsCount["status"] = status
	}

	err = r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
	if err != nil {
		r.logger.Errorf("Error getting motos count: %s", err)
		return nil, 0, err
	}

	return motors, count, nil
}

func (r *CarsPsqlRepository) GetMotoByID(ctx context.Context, id int64) (models.Moto, error) {
	var motor models.Moto

	query := `
		SELECT
			ms.id, ms.user_id, u.full_name, ms.stock_id, s.store_name, ms.brand_id, b.name,
			ms.type_motorcycles, ms.year, ms.price, ms.volume, ms.engine_type, 
			ms.number_of_clock_cycles, t.model_id, m.name, ms.air_type, ms.color, ms.vin, 
			ms.description, ms.city_id, cs.name_tm, cs.name_en, cs.name_ru,
			ms.name, ms.mail, ms.phone_number, ms.options, ms.is_comment, 
			ms.is_exchange, ms.is_credit, ms.images, ms.status
		FROM motos ms
			LEFT JOIN users u ON u.id = ms.user_id 
			LEFT JOIN stocks s ON s.id = ms.stock_id
			LEFT JOIN brands b ON b.id = ms.brand_id
			LEFT JOIN models m ON m.id = ms.model_id
			LEFT JOIN body_types bt ON bt.id = ms.body_id
			LEFT JOIN cities cs ON cs.id = ms.city_id
		WHERE ms.id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&motor.Id,
		&motor.UserId,
		&motor.UserName,
		&motor.StockId,
		&motor.StoreName,
		&motor.BrandId,
		&motor.BrandName,
		&motor.TypeMotorcycles,
		&motor.Year,
		&motor.Price,
		&motor.Volume,
		&motor.EngineType,
		&motor.NumberOfClockCycles,
		&motor.ModelId,
		&motor.ModelName,
		&motor.AirType,
		&motor.Color,
		&motor.Vin,
		&motor.Description,
		&motor.CityId,
		&motor.CityNameTM,
		&motor.CityNameEN,
		&motor.CityNameRU,
		&motor.Name,
		&motor.Mail,
		&motor.PhoneNumber,
		&motor.Options,
		&motor.IsComment,
		&motor.IsExchange,
		&motor.IsCredit,
		&motor.Images,
		&motor.Status,
	)

	if err != nil {
		r.logger.Errorf("Error getting truck by id: %s", err)
		return motor, err
	}

	return motor, nil
}

func (r *CarsPsqlRepository) GetUserByMotoId(ctx context.Context, motoId int64) (int64, error) {
	var userId int64

	query := `
		SELECT 
			user_id
		FROM motos
		WHERE id = @moto_id
	`

	args := pgx.NamedArgs{
		"moto_id": motoId,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&userId)
	if err != nil {
		r.logger.Errorf("get user by moto id err: %v", err)
		return userId, err
	}
	return userId, nil
}

func (r *CarsPsqlRepository) UpdateMotoStatus(ctx context.Context, id int64, status string) (int64, error) {
	var motoId int64

	query := `
   		UPDATE motos SET
   		     status = @status
   		WHERE id = @id
   		RETURNING id
	`

	args := pgx.NamedArgs{
		"status": status,
		"id":     id,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&motoId)
	if err != nil {
		r.logger.Errorf("update moto status err: %v", err)
		return motoId, err
	}

	return motoId, nil
}
