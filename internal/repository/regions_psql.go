package repository

import (
	"autotm-admin/internal/models"
	"context"
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

	query := `INSERT INTO regions (name_tm, name_en, name_ru) VALUES ($1, $2, $3) RETURNING id`

	err := r.client.QueryRow(ctx, query, region.NameTM, region.NameEN, region.NameRU).Scan(&id)
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
		WHERE (name_tm ILIKE '%' || $1 || '%' OR name_ru ILIKE '%' || $1 || '%' OR name_en ILIKE '%' || $1 || '%')
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.client.Query(ctx, query, search, limit, page)
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
			WHERE (name_tm ILIKE '%' || $1 || '%' OR name_ru ILIKE '%' || $1 || '%' OR name_en ILIKE '%' || $1 || '%')
		`
	errCount := r.client.QueryRow(ctx, queryCount, search).Scan(&count)
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
		    name_tm = $1, name_ru = $2, name_en = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id
	`
	err := r.client.QueryRow(ctx, query, region.NameTM, region.NameRU, region.NameEN, region.ID).Scan(&id)
	if err != nil {
		r.logger.Errorf("update region err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *RegionsPsqlRepository) DeleteRegion(ctx context.Context, id models.ID) error {
	query := `DELETE FROM regions WHERE id = $1`
	_, err := r.client.Exec(ctx, query, id.ID)
	if err != nil {
		r.logger.Errorf("delete region err: %v", err)
		return err
	}
	return nil
}

// Cities
func (r *RegionsPsqlRepository) CreateCity(ctx context.Context, city models.City) (int64, error) {
	var id int64

	query := `INSERT INTO cities (name_tm, name_en, name_ru, region_id) VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.client.QueryRow(ctx, query, city.NameTM, city.NameEN, city.NameRU, city.RegionID).Scan(&id)
	if err != nil {
		r.logger.Errorf("create city err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *RegionsPsqlRepository) GetAllCities(ctx context.Context, limit, page int64, search string) ([]models.City, int64, error) {
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
		WHERE (c.name_tm ILIKE '%' || $1 || '%' OR c.name_ru ILIKE '%' || $1 || '%' OR c.name_en ILIKE '%' || $1 || '%' 
		    OR r.name_tm ILIKE '%' || $1 || '%' OR r.name_ru ILIKE '%' || $1 || '%' OR r.name_en ILIKE '%' || $1 || '%')
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.client.Query(ctx, query, search, limit, page)
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
		WHERE (c.name_tm ILIKE '%' || $1 || '%' OR c.name_ru ILIKE '%' || $1 || '%' OR c.name_en ILIKE '%' || $1 || '%' 
		    OR r.name_tm ILIKE '%' || $1 || '%' OR r.name_ru ILIKE '%' || $1 || '%' OR r.name_en ILIKE '%' || $1 || '%')
		`
	errCount := r.client.QueryRow(ctx, queryCount, search).Scan(&count)
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
		    name_tm = $1, name_ru = $2, name_en = $3, region_id = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id
	`
	err := r.client.QueryRow(ctx, query, city.NameTM, city.NameRU, city.NameEN, city.RegionID, city.ID).Scan(&id)
	if err != nil {
		r.logger.Errorf("update city err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *RegionsPsqlRepository) DeleteCity(ctx context.Context, id models.ID) error {
	query := `DELETE FROM cities WHERE id = $1`
	_, err := r.client.Exec(ctx, query, id.ID)
	if err != nil {
		r.logger.Errorf("delete cities err: %v", err)
		return err
	}
	return nil
}
