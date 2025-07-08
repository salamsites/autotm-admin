package repository

import (
	"autotm-admin/internal/models"
	"context"
	slog "github.com/salamsites/package-log"
	spsql "github.com/salamsites/package-psql"
)

type BrandPsqlRepository struct {
	logger *slog.Logger
	client spsql.Client
}

func NewBrandPsqlRepository(logger *slog.Logger, client spsql.Client) *BrandPsqlRepository {
	return &BrandPsqlRepository{
		logger: logger,
		client: client,
	}
}

func (r *BrandPsqlRepository) CreateBrand(ctx context.Context, brand models.Brand) (int64, error) {
	var id int64

	query := `INSERT INTO brands (name, logo_path) VALUES ($1, $2) RETURNING id`

	err := r.client.QueryRow(ctx, query, brand.Name, brand.LogoPath).Scan(&id)
	if err != nil {
		r.logger.Errorf("create err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *BrandPsqlRepository) GetBrands(ctx context.Context, limit, page int64, search string) ([]models.Brand, int64, error) {
	var (
		brands []models.Brand
		count  int64
	)

	query := `
		SELECT 
		    id, name, logo_path 
		FROM brands
		WHERE name ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.client.Query(ctx, query, search, limit, page)
	if err != nil {
		r.logger.Errorf("get brands query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var brand models.Brand
		if err := rows.Scan(&brand.ID, &brand.Name, &brand.LogoPath); err != nil {
			r.logger.Errorf("get brands scan err : %v", err)
			return nil, 0, err
		}
		brands = append(brands, brand)
	}

	queryCount := `
			SELECT 
			    COUNT(*) 
			FROM brands
			WHERE name ILIKE '%' || $1 || '%'
		`
	errCount := r.client.QueryRow(ctx, queryCount, search).Scan(&count)
	if errCount != nil {
		r.logger.Errorf("get brands count err : %v", err)
		return nil, 0, err
	}
	return brands, count, nil
}

func (r *BrandPsqlRepository) UpdateBrand(ctx context.Context, brand models.Brand) (int64, error) {
	var id int64

	query := `
		UPDATE brands SET 
		    name = $1, logo_path = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id
	`
	err := r.client.QueryRow(ctx, query, brand.Name, brand.LogoPath, brand.ID).Scan(&id)
	if err != nil {
		r.logger.Errorf("update brand err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *BrandPsqlRepository) GetBrandByID(ctx context.Context, id int64) (models.Brand, error) {
	var brand models.Brand

	query := `
		SELECT
			id, name, logo_path
		FROM brands
		WHERE id = $1
	`
	err := r.client.QueryRow(ctx, query, id).Scan(&brand.ID, &brand.Name, &brand.LogoPath)
	if err != nil {
		r.logger.Errorf("get brand by id query err : %v", err)
		return brand, err
	}
	return brand, nil
}

func (r *BrandPsqlRepository) DeleteBrand(ctx context.Context, id models.ID) error {
	query := `DELETE FROM brands WHERE id = $1`
	_, err := r.client.Exec(ctx, query, id.ID)
	if err != nil {
		r.logger.Errorf("delete brand err: %v", err)
		return err
	}
	return nil
}

func (r *BrandPsqlRepository) CreateBrandModel(ctx context.Context, model models.BrandModel) (int64, error) {
	var id int64

	query := ` INSERT INTO brand_models (name, brand_id) VALUES ($1, $2) RETURNING id `

	err := r.client.QueryRow(ctx, query, model.Name, model.BrandID).Scan(&id)
	if err != nil {
		r.logger.Errorf("create brand model err : %v", err)
		return id, err
	}
	return id, nil
}

func (r *BrandPsqlRepository) GetBrandModels(ctx context.Context, limit, page int64, search string) ([]models.BrandModel, int64, error) {
	var (
		brandModels []models.BrandModel
		count       int64
	)

	query := `
		SELECT 
		    bm.id, bm.name, b.logo_path, bm.brand_id, b.name 
		FROM brand_models bm
			LEFT JOIN brands b ON bm.brand_id = b.id
		WHERE ( bm.name ILIKE '%' || $1 || '%' OR b.name ILIKE '%' || $1 || '%' )
		ORDER BY bm.created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.client.Query(ctx, query, search, limit, page)
	if err != nil {
		r.logger.Errorf("get brand models query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var brandModel models.BrandModel
		if err := rows.Scan(&brandModel.ID, &brandModel.Name, &brandModel.LogoPath, &brandModel.BrandID, &brandModel.BrandName); err != nil {
			r.logger.Errorf("get brand models scan err : %v", err)
			return nil, 0, err
		}
		brandModels = append(brandModels, brandModel)
	}

	queryCount := `
		SELECT 
			COUNT(bm.id) 
		FROM brand_models bm
		LEFT JOIN brands b ON bm.brand_id = b.id
		WHERE ( bm.name ILIKE '%' || $1 || '%' OR b.name ILIKE '%' || $1 || '%' )
		`
	errCount := r.client.QueryRow(ctx, queryCount, search).Scan(&count)
	if errCount != nil {
		r.logger.Errorf("get brand models count err : %v", err)
		return nil, 0, err
	}
	return brandModels, count, nil
}

func (r *BrandPsqlRepository) UpdateBrandModel(ctx context.Context, model models.BrandModel) (int64, error) {
	var id int64

	query := `
		UPDATE brand_models SET 
		    name = $1, brand_id = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id;
	`
	err := r.client.QueryRow(ctx, query, model.Name, model.BrandID, model.ID).Scan(&id)
	if err != nil {
		r.logger.Errorf("update brand model err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *BrandPsqlRepository) DeleteBrandModel(ctx context.Context, id models.ID) error {
	query := `DELETE FROM brand_models WHERE id = $1`
	_, err := r.client.Exec(ctx, query, id.ID)
	if err != nil {
		r.logger.Errorf("delete brand model err: %v", err)
		return err
	}
	return nil
}
