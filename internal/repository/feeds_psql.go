package repository

import (
	"autotm-admin/internal/models"
	"context"
	"strings"

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
			cr.id,
			-- Stock
			cr.stock_id, s.store_name,
			-- Brand
			cr.brand_id, b.name,
			-- Model
			cr.model_id, m.name, 
			-- City
			cr.city_id, cs.name_tm, cs.name_en, cs.name_ru, 
			cr.name, cr.mail, cr.phone_number, 
			cr.year, cr.price, cr.images, 
			cr.status, cr.created_at
		FROM cars cr
			LEFT JOIN stocks s ON s.id = cr.stock_id
			LEFT JOIN brands b ON b.id = cr.brand_id
			LEFT JOIN models m ON m.id = cr.model_id
			LEFT JOIN cities cs ON cs.id = cr.city_id
	`

	conditions := []string{}
	args := pgx.NamedArgs{
		"search": search,
		"limit":  limit,
		"offset": page,
	}

	// Search condition
	if search != "" {
		searchCondition := `
			(s.store_name ILIKE '%' || @search || '%' OR 
			 b.name ILIKE '%' || @search || '%' OR
			 m.name ILIKE '%' || @search || '%' OR
			 cs.name_tm ILIKE '%' || @search || '%' OR
			 cs.name_en ILIKE '%' || @search || '%' OR
			 cs.name_ru ILIKE '%' || @search || '%' OR
			 cr.name ILIKE '%' || @search || '%' OR
			 cr.phone_number ILIKE '%' || @search || '%' OR
			 CAST(cr.year AS TEXT) ILIKE '%' || @search || '%' OR
			 CAST(cr.price AS TEXT) ILIKE '%' || @search || '%')
		`
		conditions = append(conditions, searchCondition)
		args["search"] = search
	}

	if status != "" {
		conditions = append(conditions, "cr.status = @status")
		args["status"] = status
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
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
			&car.StockId,
			&car.StoreName,
			&car.BrandId,
			&car.BrandName,
			&car.ModelId,
			&car.ModelName,
			&car.CityId,
			&car.CityNameTM,
			&car.CityNameEN,
			&car.CityNameRU,
			&car.Name,
			&car.Mail,
			&car.PhoneNumber,
			&car.Year,
			&car.Price,
			&car.Images,
			&car.Status,
			&car.CreatedAt,
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
				LEFT JOIN stocks s ON s.id = cr.stock_id
				LEFT JOIN brands b ON b.id = cr.brand_id
				LEFT JOIN models m ON m.id = cr.model_id
				LEFT JOIN cities cs ON cs.id = cr.city_id
		`
	countArgs := pgx.NamedArgs{}
	if search != "" {
		countArgs["search"] = search
	}
	if status != "" {
		countArgs["status"] = status
	}

	if len(conditions) > 0 {
		queryCount += " WHERE " + strings.Join(conditions, " AND ")
	}

	err = r.client.QueryRow(ctx, queryCount, countArgs).Scan(&count)
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
			cr.id,
			-- User
			cr.user_id, u.full_name, 
			-- Stock
			cr.stock_id, s.store_name, 
			-- Brand
			cr.brand_id, b.name,
			-- Model
			cr.model_id, m.name, 
			cr.year, cr.mileage, cr.color, cr.engine_capacity, cr.engine_type,
			-- Body
			cr.body_id, bt.name_tm, bt.name_en, bt.name_ru, 
			cr.transmission, cr.drive_type, cr.vin, cr.description, 
			-- City
			cr.city_id, cs.name_tm, cs.name_en, cs.name_ru, 
			cr.name, cr.mail, cr.phone_number, cr.price, cr.is_comment, 
			cr.is_exchange, cr.is_credit, cr.images, cr.status, cr.created_at, cr.updated_at
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
		&car.BrandId, &car.BrandName, &car.ModelId, &car.ModelName, &car.Year, &car.Mileage, &car.Color,
		&car.EngineCapacity, &car.EngineType, &car.BodyId, &car.BodyNameTM, &car.BodyNameEN, &car.BodyNameRU, &car.Transmission,
		&car.DriveType, &car.Vin, &car.Description, &car.CityId, &car.CityNameTM, &car.CityNameEN, &car.CityNameRU, &car.Name, &car.Mail,
		&car.PhoneNumber, &car.Price, &car.IsComment, &car.IsExchange, &car.IsCredit, &car.Images, &car.Status, &car.CreatedAt, &car.UpdatedAt,
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

func (r *CarsPsqlRepository) GetTrucks(ctx context.Context, limit, page int64, search, status string) ([]models.Truck, int64, error) {
	var (
		trucks []models.Truck
		count  int64
	)

	query := `
		SELECT
			t.id, 
			-- Stock
			t.stock_id, s.store_name, 
			-- Brand
			t.brand_id, b.name,
			-- Model
			t.model_id, m.name,
			-- City
			t.city_id, cs.name_tm, cs.name_en, cs.name_ru, 
			t.name, t.mail, t.phone_number, t.year,
			t.price, t.images, t.status, t.created_at
		FROM trucks t
			LEFT JOIN stocks s ON s.id = t.stock_id
			LEFT JOIN brands b ON b.id = t.brand_id
			LEFT JOIN models m ON m.id = t.model_id
			LEFT JOIN cities cs ON cs.id = t.city_id
	`

	conditions := []string{}
	args := pgx.NamedArgs{
		"search": search,
		"limit":  limit,
		"offset": page,
	}

	// Search condition
	if search != "" {
		searchCondition := `
			(s.store_name ILIKE '%' || @search || '%' OR 
			 b.name ILIKE '%' || @search || '%' OR
			 m.name ILIKE '%' || @search || '%' OR
			 cs.name_tm ILIKE '%' || @search || '%' OR
			 cs.name_en ILIKE '%' || @search || '%' OR
			 cs.name_ru ILIKE '%' || @search || '%' OR
			 t.name ILIKE '%' || @search || '%' OR
			 t.phone_number ILIKE '%' || @search || '%' OR
			 CAST(t.year AS TEXT) ILIKE '%' || @search || '%' OR
			 CAST(t.price AS TEXT) ILIKE '%' || @search || '%')
		`
		conditions = append(conditions, searchCondition)
		args["search"] = search
	}

	if status != "" {
		conditions = append(conditions, "t.status = @status")
		args["status"] = status
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
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
			&truck.StockId,
			&truck.StoreName,
			&truck.BrandId,
			&truck.BrandName,
			&truck.Price,
			&truck.ModelId,
			&truck.ModelName,
			&truck.Year,
			&truck.CityId,
			&truck.CityNameTM,
			&truck.CityNameEN,
			&truck.CityNameRU,
			&truck.Name,
			&truck.Mail,
			&truck.PhoneNumber,
			&truck.Year,
			&truck.Price,
			&truck.Images,
			&truck.Status,
			&truck.CreatedAt,
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
				LEFT JOIN stocks s ON s.id = t.stock_id
				LEFT JOIN brands b ON b.id = t.brand_id
				LEFT JOIN models m ON m.id = t.model_id
				LEFT JOIN cities cs ON cs.id = t.city_id
		`
	countArgs := pgx.NamedArgs{}
	if search != "" {
		countArgs["search"] = search
	}
	if status != "" {
		countArgs["status"] = status
	}

	if len(conditions) > 0 {
		queryCount += " WHERE " + strings.Join(conditions, " AND ")
	}

	err = r.client.QueryRow(ctx, queryCount, countArgs).Scan(&count)
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
			t.id, 
			-- User
			t.user_id, u.full_name, 
			-- Stock
			t.stock_id, s.store_name, 
			-- Brand
			t.brand_id, b.name,
			t.load_capacity, t.price, t.body_type, t.drive_type, t.transmission, t.engine_type,
			-- Model
			t.model_id, m.name, 
			t.year, t.seats, t.cab_type, t.wheel_formula, t.chassis, t.cab_suspension,
			t.bus_type, t.suspension_type, t.brakes, t.axles, t.engine_hours, t.vehicle_type, t.engine_capacity,
			t.forklift_type, t.lifting_capacity, t.mileage, t.excavator_type, t.bulldozer_type, t.color, t.vin, 
			-- Body
			t.body_id, bt.name_tm, bt.name_en, bt.name_ru, 
			t.description,
			-- City
			t.city_id, cs.name_tm, cs.name_en, cs.name_ru, 
			t.name, t.mail, t.phone_number, t.is_comment, t.is_exchange, 
			t.is_credit, t.images, t.status, t.created_at, t.updated_at
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
		&truck.IsComment, &truck.IsExchange, &truck.IsCredit, &truck.Images, &truck.Status, &truck.CreatedAt, &truck.UpdatedAt,
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

func (r *CarsPsqlRepository) GetMotors(ctx context.Context, limit, page int64, search, status string) ([]models.Moto, int64, error) {
	var (
		motors []models.Moto
		count  int64
	)

	query := `
		SELECT
			ms.id, 
			-- Stock
			ms.stock_id, s.store_name, 
			-- Brand
			ms.brand_id, b.name,
			-- Model
			ms.model_id, m.name, 
			-- City
			ms.city_id, cs.name_tm, cs.name_en, cs.name_ru,
			ms.name, ms.mail, ms.phone_number, 
			ms.year, ms.price,
			ms.images, ms.status, ms.created_at
		FROM motoes ms
			LEFT JOIN stocks s ON s.id = ms.stock_id
			LEFT JOIN brands b ON b.id = ms.brand_id
			LEFT JOIN models m ON m.id = ms.model_id
			LEFT JOIN cities cs ON cs.id = ms.city_id
	`

	conditions := []string{}
	args := pgx.NamedArgs{
		"search": search,
		"limit":  limit,
		"offset": page,
	}

	// Search condition
	if search != "" {
		searchCondition := `
			(s.store_name ILIKE '%' || @search || '%' OR 
			 b.name ILIKE '%' || @search || '%' OR
			 m.name ILIKE '%' || @search || '%' OR
			 cs.name_tm ILIKE '%' || @search || '%' OR
			 cs.name_en ILIKE '%' || @search || '%' OR
			 cs.name_ru ILIKE '%' || @search || '%' OR
			 ms.name ILIKE '%' || @search || '%' OR
			 ms.phone_number ILIKE '%' || @search || '%' OR
			 CAST(ms.year AS TEXT) ILIKE '%' || @search || '%' OR
			 CAST(ms.price AS TEXT) ILIKE '%' || @search || '%')
		`
		conditions = append(conditions, searchCondition)
		args["search"] = search
	}

	if status != "" {
		conditions = append(conditions, "ms.status = @status")
		args["status"] = status
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
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
			&motor.StockId,
			&motor.StoreName,
			&motor.BrandId,
			&motor.BrandName,
			&motor.ModelId,
			&motor.ModelName,
			&motor.CityId,
			&motor.CityNameTM,
			&motor.CityNameEN,
			&motor.CityNameRU,
			&motor.Name,
			&motor.Mail,
			&motor.PhoneNumber,
			&motor.Year,
			&motor.Price,
			&motor.Images,
			&motor.Status,
			&motor.CreatedAt,
		)
		if err != nil {
			r.logger.Errorf("Error getting motors: %s", err)
		}
		motors = append(motors, motor)
	}

	queryCount := `
			SELECT 
    			COUNT(ms.id) 
			FROM motoes ms
				LEFT JOIN stocks s ON s.id = ms.stock_id
				LEFT JOIN brands b ON b.id = ms.brand_id
				LEFT JOIN models m ON m.id = ms.model_id
				LEFT JOIN cities cs ON cs.id = ms.city_id
		`
	countArgs := pgx.NamedArgs{}
	if search != "" {
		countArgs["search"] = search
	}
	if status != "" {
		countArgs["status"] = status
	}

	if len(conditions) > 0 {
		queryCount += " WHERE " + strings.Join(conditions, " AND ")
	}

	err = r.client.QueryRow(ctx, queryCount, countArgs).Scan(&count)
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
			ms.id, 
			-- User
			ms.user_id, u.full_name, 
			-- Stock
			ms.stock_id, s.store_name,
			-- Brand
			ms.brand_id, b.name,
			ms.type_motorcycles, ms.year, ms.price, ms.volume, ms.engine_type, 
			ms.number_of_clock_cycles, 
			-- Model
			ms.model_id, m.name, 
			ms.air_type, ms.color, ms.vin, ms.description,
			-- City
			ms.city_id, cs.name_tm, cs.name_en, cs.name_ru,
			-- Body
			ms.body_id, bt.name_tm, bt.name_en, bt.name_ru,
			ms.name, ms.mail, ms.phone_number, ms.options, ms.is_comment, 
			ms.is_exchange, ms.is_credit, ms.images, ms.status,
			ms.options, ms.created_at, ms.updated_at
		FROM motoes ms
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
		&motor.BodyId,
		&motor.BodyNameTM,
		&motor.BodyNameEN,
		&motor.BodyNameRU,
		&motor.Name,
		&motor.Mail,
		&motor.PhoneNumber,
		&motor.Options,
		&motor.IsComment,
		&motor.IsExchange,
		&motor.IsCredit,
		&motor.Images,
		&motor.Status,
		&motor.Options,
		&motor.CreatedAt,
		&motor.UpdatedAt,
	)

	if err != nil {
		r.logger.Errorf("Error getting moto by id: %s", err)
		return motor, err
	}

	return motor, nil
}

func (r *CarsPsqlRepository) UpdateMotoStatus(ctx context.Context, id int64, status string) (int64, error) {
	var motoId int64

	query := `
   		UPDATE motoes SET
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
