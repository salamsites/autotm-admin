package repository

import (
	"autotm-admin/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	slog "github.com/salamsites/package-log"
	spsql "github.com/salamsites/package-psql"
)

type RegionsPsqlRepository struct {
	logger *slog.Logger
	client spsql.Client
}

func NewRegionsPsqlRepository(logger *slog.Logger, client spsql.Client) *RegionsPsqlRepository {
	return &RegionsPsqlRepository{
		logger: logger,
		client: client,
	}
}

func (r *RegionsPsqlRepository) CreateRegion(ctx context.Context, region models.Region) (int64, error) {
	var id int64

	query := `INSERT INTO regions (name_tm, name_en, name_ru) VALUES (@name_tm, @name_en, @name_ru) RETURNING id`

	args := pgx.NamedArgs{
		"name_tm": region.NameTM,
		"name_en": region.NameEN,
		"name_ru": region.NameRU,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("create region err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *RegionsPsqlRepository) GetAllRegions(ctx context.Context, limit, page int64, search string) ([]models.Region, int64, error) {
	var (
		regions []models.Region
		count   int64
	)

	query := `
		SELECT 
		    id, name_tm, name_en, name_ru 
		FROM regions
		WHERE (name_tm ILIKE '%' || @search || '%' OR name_ru ILIKE '%' || @search || '%' OR name_en ILIKE '%' || @search || '%')
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	args := pgx.NamedArgs{
		"search": search,
		"limit":  limit,
		"offset": page,
	}
	rows, err := r.client.Query(ctx, query, args)
	if err != nil {
		r.logger.Errorf("get all regions query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var region models.Region
		if err := rows.Scan(&region.ID, &region.NameTM, &region.NameEN, &region.NameRU); err != nil {
			r.logger.Errorf("get all regions scan err : %v", err)
			return nil, 0, err
		}
		regions = append(regions, region)
	}

	queryCount := `
			SELECT 
			    COUNT(*) 
			FROM regions
			WHERE (name_tm ILIKE '%' || @search || '%' OR name_ru ILIKE '%' || @search || '%' OR name_en ILIKE '%' || @search || '%')
		`

	argsCount := pgx.NamedArgs{
		"search": search,
	}
	errCount := r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
	if errCount != nil {
		r.logger.Errorf("get all regions count err : %v", err)
		return nil, 0, err
	}
	return regions, count, nil
}

func (r *RegionsPsqlRepository) UpdateRegion(ctx context.Context, region models.Region) (int64, error) {
	var id int64

	query := `
		UPDATE regions SET 
		    name_tm = @name_tm, name_ru = @name_ru, name_en = @name_en, updated_at = NOW()
		WHERE id = @id
		RETURNING id
	`
	args := pgx.NamedArgs{
		"name_tm": region.NameTM,
		"name_ru": region.NameRU,
		"name_en": region.NameEN,
		"id":      region.ID,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("update region err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *RegionsPsqlRepository) DeleteRegion(ctx context.Context, id int64) error {
	query := `DELETE FROM regions WHERE id = @id`

	args := pgx.NamedArgs{
		"id": id,
	}
	_, err := r.client.Exec(ctx, query, args)
	if err != nil {
		r.logger.Errorf("delete region err: %v", err)
		return err
	}
	return nil
}

// Cities
func (r *RegionsPsqlRepository) CreateCity(ctx context.Context, city models.City) (int64, error) {
	var id int64

	query := `INSERT INTO cities (name_tm, name_en, name_ru, region_id) VALUES (@name_tm, @name_en, @name_ru, @region_id) RETURNING id`

	args := pgx.NamedArgs{
		"name_tm":   city.NameTM,
		"name_en":   city.NameEN,
		"name_ru":   city.NameRU,
		"region_id": city.RegionID,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("create city err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *RegionsPsqlRepository) GetAllCities(ctx context.Context, limit, page int64, search string, regionIds []int64) ([]models.City, int64, error) {
	var (
		cities []models.City
		count  int64
	)

	query := `
		SELECT 
		    c.id, c.name_tm, c.name_en, c.name_ru, c.region_id,
		    r.name_tm, r.name_en, r.name_ru
		FROM cities c
			LEFT JOIN regions r on r.id = c.region_id
		WHERE (c.name_tm ILIKE '%' || @search || '%' OR c.name_ru ILIKE '%' || @search || '%' OR c.name_en ILIKE '%' || @search || '%' 
		    OR r.name_tm ILIKE '%' || @search || '%' OR r.name_ru ILIKE '%' || @search || '%' OR r.name_en ILIKE '%' || @search || '%')
	`

	args := pgx.NamedArgs{
		"search": search,
		"limit":  limit,
		"offset": page,
	}

	if len(regionIds) > 0 {
		query += " AND c.region_id = ANY(@regionIds)"
		args["regionIds"] = regionIds
	}

	query += `
        ORDER BY c.created_at DESC
        LIMIT @limit OFFSET @offset;
    `
	rows, err := r.client.Query(ctx, query, args)
	if err != nil {
		r.logger.Errorf("get all cities query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var city models.City
		err = rows.Scan(&city.ID, &city.NameTM, &city.NameEN, &city.NameRU, &city.RegionID,
			&city.RegionNameTM, &city.RegionNameEN, &city.RegionNameRU,
		)
		if err != nil {
			r.logger.Errorf("get all cities scan err : %v", err)
			return nil, 0, err
		}
		cities = append(cities, city)
	}

	queryCount := `
			SELECT 
			    COUNT(c.id) 
			FROM cities c
			LEFT JOIN regions r on r.id = c.region_id
		WHERE (c.name_tm ILIKE '%' || @search || '%' OR c.name_ru ILIKE '%' || @search || '%' OR c.name_en ILIKE '%' || @search || '%' 
		    OR r.name_tm ILIKE '%' || @search || '%' OR r.name_ru ILIKE '%' || @search || '%' OR r.name_en ILIKE '%' || @search || '%')
		`

	countArgs := pgx.NamedArgs{
		"search": search,
	}
	if len(regionIds) > 0 {
		queryCount += " AND c.region_id = ANY(@regionIds)"
		countArgs["regionIds"] = regionIds
	}

	errCount := r.client.QueryRow(ctx, queryCount, countArgs).Scan(&count)
	if errCount != nil {
		r.logger.Errorf("get all cities count err : %v", err)
		return nil, 0, err
	}
	return cities, count, nil
}

func (r *RegionsPsqlRepository) UpdateCity(ctx context.Context, city models.City) (int64, error) {
	var id int64

	query := `
		UPDATE cities SET 
		    name_tm = @name_tm, name_ru = @name_ru, name_en = @name_en, region_id = @region_id, updated_at = NOW()
		WHERE id = @id
		RETURNING id
	`

	args := pgx.NamedArgs{
		"name_tm":   city.NameTM,
		"name_ru":   city.NameRU,
		"name_en":   city.NameEN,
		"region_id": city.RegionID,
		"id":        city.ID,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("update city err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *RegionsPsqlRepository) DeleteCity(ctx context.Context, id int64) error {
	query := `DELETE FROM cities WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}
	_, err := r.client.Exec(ctx, query, args)
	if err != nil {
		r.logger.Errorf("delete cities err: %v", err)
		return err
	}
	return nil
}
