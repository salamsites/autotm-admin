package repository

import (
	"autotm-admin/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
	slog "github.com/salamsites/package-log"
	spsql "github.com/salamsites/package-psql"
)

type AutoStorePsqlRepository struct {
	logger *slog.Logger
	client spsql.Client
}

func NewAutoStorePsqlRepository(logger *slog.Logger, client spsql.Client) *AutoStorePsqlRepository {
	return &AutoStorePsqlRepository{
		logger: logger,
		client: client,
	}
}

func (r *AutoStorePsqlRepository) CreateAutoStore(ctx context.Context, autoStore models.AutoStore) (int64, error) {
	var id int64

	query := ` 
			INSERT INTO auto_stores 
			    (user_id, phone_number, email, store_name, images, logo_path, address, region_id, city_id) 
			VALUES (@user_id, @phone_number, @email, @store_name, @images, @logo_path, @address, @region_id, @city_id) 
			RETURNING id;
	`

	args := pgx.NamedArgs{
		"user_id":      autoStore.UserID,
		"phone_number": autoStore.PhoneNumber,
		"email":        autoStore.Email,
		"store_name":   autoStore.StoreName,
		"images":       autoStore.Images,
		"logo_path":    autoStore.LogoPath,
		"address":      autoStore.Address,
		"region_id":    autoStore.RegionID,
		"city_id":      autoStore.CityID,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("Error creating auto store: %s", err.Error())
		return id, err
	}
	return id, nil
}

func (r *AutoStorePsqlRepository) GetAutoStores(ctx context.Context, limit, page int64, search string) ([]models.AutoStore, int64, error) {
	var (
		autoStores []models.AutoStore
		count      int64
	)

	query := `
			SELECT
				ast.id, ast.user_id, ast.phone_number, ast.email, ast.store_name, 
				ast.images, ast.logo_path, ast.address, ast.city_id, c.name_tm,
				c.name_en, c.name_ru, ast.region_id, r.name_tm, r.name_en, r.name_ru
           FROM auto_stores ast
           LEFT JOIN cities c ON c.id = ast.city_id
           LEFT JOIN regions r on r.id = ast.region_id
		   WHERE ast.store_name ILIKE '%' || @search || '%'
		   ORDER BY ast.created_at DESC
		   LIMIT @limit OFFSET @page
		`

	args := pgx.NamedArgs{
		"search": search,
		"limit":  limit,
		"page":   page,
	}

	rows, err := r.client.Query(ctx, query, args)
	if err != nil {
		r.logger.Errorf("Error getting auto-store: %s", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var store models.AutoStore
		err = rows.Scan(
			&store.ID,
			&store.UserID,
			&store.PhoneNumber,
			&store.Email,
			&store.StoreName,
			&store.Images,
			&store.LogoPath,
			&store.Address,
			&store.CityID,
			&store.CityNameTM,
			&store.CityNameEN,
			&store.CityNameRU,
			&store.RegionID,
			&store.RegionNameTM,
			&store.RegionNameEN,
			&store.RegionNameRU,
		)
		if err != nil {
			r.logger.Errorf("Error scanning auto-store: %s", err)
			return nil, 0, err
		}
		autoStores = append(autoStores, store)
	}

	queryCount := `SELECT COUNT(*) FROM auto_stores WHERE store_name ILIKE '%' || @search || '%'`
	argsCount := pgx.NamedArgs{
		"search": search,
	}

	err = r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
	if err != nil {
		r.logger.Errorf("Error getting auto-store count: %s", err)
		return nil, 0, err
	}

	return autoStores, count, nil
}

func (r *AutoStorePsqlRepository) UpdateAutoStore(ctx context.Context, autoStore models.AutoStore) (int64, error) {
	var autoStoreID int64

	query := `
		UPDATE auto_stores SET 
		    user_id = @user_id, phone_number = @phone_number, email = @email, store_name = @store_name, images = @images,
		    logo_path = @logo_path, address = @address, region_id = @region_id, city_id = @city_id, updated_at = NOW()
		WHERE id = @id
		RETURNING id
	`

	args := pgx.NamedArgs{
		"user_id":      autoStore.UserID,
		"phone_number": autoStore.PhoneNumber,
		"email":        autoStore.Email,
		"store_name":   autoStore.StoreName,
		"images":       autoStore.Images,
		"logo_path":    autoStore.LogoPath,
		"address":      autoStore.Address,
		"region_id":    autoStore.RegionID,
		"city_id":      autoStore.CityID,
		"id":           autoStore.ID,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&autoStoreID)
	if err != nil {
		r.logger.Errorf("update body types err: %v", err)
		return autoStoreID, err
	}
	return autoStoreID, nil
}

func (r *AutoStorePsqlRepository) DeleteAutoStore(ctx context.Context, id models.ID) error {
	query := ` DELETE FROM auto_stores WHERE id = @id `

	args := pgx.NamedArgs{
		"id": id.ID,
	}
	_, err := r.client.Exec(ctx, query, args)
	if err != nil {
		r.logger.Errorf("delete autoStore err: %v", err)
		return err
	}
	return nil
}
