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
			    (user_id, phone_number, email, store_name, images, logo, address, region_id, city_id) 
			VALUES (@user_id, @phone_number, @email, @store_name, @images, @logo, @address, @region_id, @city_id) 
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
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("Error creating stock: %s", err.Error())
		return id, err
	}
	return id, nil
}

func (r *StockPsqlRepository) GetStocks(ctx context.Context, limit, page int64, search string) ([]models.Stock, int64, error) {
	var (
		stocks []models.Stock
		count  int64
	)

	query := `
			SELECT
				s.id, s.user_id, u.full_name, s.phone_number, s.email, s.store_name, 
				s.images, s.logo, s.address, s.city_id, c.name_tm,
				c.name_en, c.name_ru, s.region_id, r.name_tm, r.name_en, r.name_ru
           FROM stocks s
           LEFT JOIN users u ON u.id = s.user_id
           LEFT JOIN cities c ON c.id = s.city_id
           LEFT JOIN regions r on r.id = s.region_id
		   WHERE s.store_name ILIKE '%' || @search || '%' OR u.full_name ILIKE '%' || @search || '%'
		   ORDER BY s.created_at DESC
		   LIMIT @limit OFFSET @page
		`

	args := pgx.NamedArgs{
		"search": search,
		"limit":  limit,
		"page":   page,
	}

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
			WHERE s.store_name ILIKE '%' || @search || '%' OR u.full_name ILIKE '%' || @search || '%'
		`
	argsCount := pgx.NamedArgs{
		"search": search,
	}

	err = r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
	if err != nil {
		r.logger.Errorf("Error getting auto-store count: %s", err)
		return nil, 0, err
	}

	return stocks, count, nil
}

func (r *StockPsqlRepository) UpdateStock(ctx context.Context, stock models.Stock) (int64, error) {
	var stockID int64

	query := `
		UPDATE stocks SET 
		    user_id = @user_id, phone_number = @phone_number, email = @email, store_name = @store_name, 
		    images = @images, logo = @logo, address = @address, region_id = @region_id, city_id = @city_id
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
		"id":           stock.ID,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&stockID)
	if err != nil {
		r.logger.Errorf("update stocks err: %v", err)
		return stockID, err
	}
	return stockID, nil
}

func (r *StockPsqlRepository) DeleteStock(ctx context.Context, id models.ID) error {
	query := ` DELETE FROM stocks WHERE id = @id `

	args := pgx.NamedArgs{
		"id": id.ID,
	}
	_, err := r.client.Exec(ctx, query, args)
	if err != nil {
		r.logger.Errorf("delete stock err: %v", err)
		return err
	}
	return nil
}
