package repository

import (
	"autotm-admin/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	slog "github.com/salamsites/package-log"
	spsql "github.com/salamsites/package-psql"
)

type StockPsqlRepository struct {
	logger *slog.Logger
	client spsql.Client
}

func NewStockPsqlRepository(logger *slog.Logger, client spsql.Client) *StockPsqlRepository {
	return &StockPsqlRepository{
		logger: logger,
		client: client,
	}
}

func (r *StockPsqlRepository) CreateStock(ctx context.Context, stock models.Stock) (int64, error) {
	var id int64

	query := ` 
			INSERT INTO stocks 
			    (user_id, phone_number, email, store_name, images, logo, address, region_id, city_id, status, description) 
			VALUES (@user_id, @phone_number, @email, @store_name, @images, @logo, @address, @region_id, @city_id, @status, @description) 
			RETURNING id;
	`

	args := pgx.NamedArgs{
		"user_id":      stock.UserID,
		"phone_number": stock.PhoneNumber,
		"email":        stock.Email,
		"store_name":   stock.StoreName,
		"images":       stock.Images,
		"logo":         stock.Logo,
		"address":      stock.Address,
		"region_id":    stock.RegionID,
		"city_id":      stock.CityID,
		"status":       stock.Status,
		"description":  stock.Description,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("Error creating stock: %s", err.Error())
		return id, err
	}
	return id, nil
}

func (r *StockPsqlRepository) UpdateStockImages(ctx context.Context, stockID int64, images []string) error {
	query := `UPDATE stocks SET images = @images WHERE id = @id`

	args := pgx.NamedArgs{
		"images": images,
		"id":     stockID,
	}
	_, err := r.client.Exec(ctx, query, args)
	if err != nil {
		r.logger.Errorf("update images err: %v", err)
		return err
	}
	return nil
}

func (r *StockPsqlRepository) UpdateStockLogo(ctx context.Context, stockID int64, logo string) error {
	query := `UPDATE stocks SET logo = @logo WHERE id = @id`

	args := pgx.NamedArgs{
		"logo": logo,
		"id":   stockID,
	}
	_, err := r.client.Exec(ctx, query, args)
	if err != nil {
		r.logger.Errorf("update logo err: %v", err)
		return err
	}
	return nil
}

func (r *StockPsqlRepository) GetStocks(ctx context.Context, limit, page int64, search, status string) ([]models.Stock, int64, error) {
	var (
		stocks []models.Stock
		count  int64
	)

	query := `
			SELECT
				s.id, s.user_id, u.full_name, s.phone_number, s.email, s.store_name, 
				s.images, s.logo, s.address, s.city_id, c.name_tm, c.name_en, c.name_ru, 
				s.region_id, r.name_tm, r.name_en, r.name_ru, s.status, s.description
           FROM stocks s
           LEFT JOIN users u ON u.id = s.user_id
           LEFT JOIN cities c ON c.id = s.city_id
           LEFT JOIN regions r on r.id = s.region_id
		   WHERE (s.store_name ILIKE '%' || @search || '%' OR u.full_name ILIKE '%' || @search || '%')
		`

	args := pgx.NamedArgs{
		"search": search,
		"limit":  limit,
		"offset": page,
	}

	if status != "" {
		query += " AND s.status = @status "
		args["status"] = status
	}

	query += `
   		ORDER BY s.created_at DESC
   		LIMIT @limit OFFSET @offset
    `
	rows, err := r.client.Query(ctx, query, args)
	if err != nil {
		r.logger.Errorf("Error getting stock: %s", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err = rows.Scan(
			&stock.ID,
			&stock.UserID,
			&stock.UserName,
			&stock.PhoneNumber,
			&stock.Email,
			&stock.StoreName,
			&stock.Images,
			&stock.Logo,
			&stock.Address,
			&stock.CityID,
			&stock.CityNameTM,
			&stock.CityNameEN,
			&stock.CityNameRU,
			&stock.RegionID,
			&stock.RegionNameTM,
			&stock.RegionNameEN,
			&stock.RegionNameRU,
			&stock.Status,
			&stock.Description,
		)
		if err != nil {
			r.logger.Errorf("Error scanning stock: %s", err)
			return nil, 0, err
		}
		stocks = append(stocks, stock)
	}

	queryCount := `
			SELECT 
    			COUNT(s.id) 
			FROM stocks s 
				LEFT JOIN users u ON u.id = s.user_id
			    LEFT JOIN cities c ON c.id = s.city_id
			    LEFT JOIN regions r on r.id = s.region_id
			WHERE (s.store_name ILIKE '%' || @search || '%' OR u.full_name ILIKE '%' || @search || '%')
		`
	argsCount := pgx.NamedArgs{
		"search": search,
	}

	if status != "" {
		queryCount += " AND s.status = @status "
		argsCount["status"] = status
	}

	err = r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
	if err != nil {
		r.logger.Errorf("Error getting stocks count: %s", err)
		return nil, 0, err
	}

	return stocks, count, nil
}

func (r *StockPsqlRepository) GetStockByID(ctx context.Context, stockID int64) (models.Stock, error) {
	var stock models.Stock

	query := `
		SELECT
			s.id, s.user_id, u.full_name, s.phone_number, 
			s.email, s.store_name, s.images, s.logo, s.region_id, 
			r.name_tm, r.name_ru, r.name_en, s.city_id, c.name_tm,
			c.name_en, c.name_ru, s.address, s.status, s.description
		FROM stocks s
			LEFT JOIN users u ON u.id = s.user_id
			LEFT JOIN cities c ON c.id = s.city_id
			LEFT JOIN regions r ON r.id = s.region_id
		WHERE s.id = @id
	`

	args := pgx.NamedArgs{
		"id": stockID,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&stock.ID, &stock.UserID, &stock.UserName, &stock.PhoneNumber, &stock.Email,
		&stock.StoreName, &stock.Images, &stock.Logo, &stock.RegionID, &stock.RegionNameTM, &stock.RegionNameRU, &stock.RegionNameEN,
		&stock.CityID, &stock.CityNameTM, &stock.CityNameEN, &stock.CityNameRU, &stock.Address, &stock.Status, &stock.Description,
	)

	if err != nil {
		r.logger.Errorf("Error getting stock by id: %s", err)
		return stock, err
	}

	return stock, nil
}

func (r *StockPsqlRepository) UpdateStock(ctx context.Context, stock models.Stock) (int64, error) {
	var stockID int64

	query := `
		UPDATE stocks SET 
		    user_id = @user_id, phone_number = @phone_number, email = @email, store_name = @store_name, images = @images, 
		    logo = @logo, address = @address, region_id = @region_id, city_id = @city_id, status = @status, description = @description
		WHERE id = @id
		RETURNING id
	`

	args := pgx.NamedArgs{
		"user_id":      stock.UserID,
		"phone_number": stock.PhoneNumber,
		"email":        stock.Email,
		"store_name":   stock.StoreName,
		"images":       stock.Images,
		"logo":         stock.Logo,
		"address":      stock.Address,
		"region_id":    stock.RegionID,
		"city_id":      stock.CityID,
		"status":       stock.Status,
		"description":  stock.Description,
		"id":           stock.ID,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&stockID)
	if err != nil {
		r.logger.Errorf("update stocks err: %v", err)
		return stockID, err
	}
	return stockID, nil
}

func (r *StockPsqlRepository) DeleteStock(ctx context.Context, id int64) error {
	query := ` DELETE FROM stocks WHERE id = @id `

	args := pgx.NamedArgs{
		"id": id,
	}
	_, err := r.client.Exec(ctx, query, args)
	if err != nil {
		r.logger.Errorf("delete stock err: %v", err)
		return err
	}
	return nil
}

func (r *StockPsqlRepository) UpdateStockStatus(ctx context.Context, id int64, status string) (int64, error) {
	var stockID int64

	query := `
   		UPDATE stocks SET
   		     status = @status
   		WHERE id = @id
   		RETURNING id
	`

	args := pgx.NamedArgs{
		"status": status,
		"id":     id,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&stockID)
	if err != nil {
		r.logger.Errorf("update stock status err: %v", err)
		return stockID, err
	}

	return stockID, nil
}

func (r *StockPsqlRepository) GetUserByStockId(ctx context.Context, stockId int64) (int64, error) {
	var userId int64

	query := `
		SELECT 
			user_id
		FROM stocks
		WHERE stock_id = @stock_id
	`

	args := pgx.NamedArgs{
		"stock_id": stockId,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&userId)
	if err != nil {
		r.logger.Errorf("get user by stock id err: %v", err)
		return userId, err
	}
	return userId, nil
}
