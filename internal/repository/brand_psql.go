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

func (r *BrandPsqlRepository) CreateBodyType(ctx context.Context, bodyType models.BodyType) (int64, error) {
	var id int64

	query := ` INSERT INTO body_types (name, image_path) VALUES ($1, $2) RETURNING id `

	err := r.client.QueryRow(ctx, query, bodyType.Name, bodyType.ImagePath).Scan(&id)
	if err != nil {
		r.logger.Errorf("Error creating body type: %s", err.Error())
		return id, err
	}
	return id, nil
}

func (r *BrandPsqlRepository) GetBodyType(ctx context.Context, limit, page int64, search string) ([]models.BodyType, int64, error) {
	var (
		bodyTypes []models.BodyType
		count     int64
	)

	query := `
			SELECT 
				id, name, category, image_path
            FROM body_types
			WHERE name ILIKE '%' || $2 || '%'
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4;
		`

	rows, err := r.client.Query(ctx, query, search, limit, page)
	if err != nil {
		r.logger.Errorf("get body types query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var bodyType models.BodyType
		if err = rows.Scan(&bodyType.ID, &bodyType.Name, &bodyType.ImagePath); err != nil {
			r.logger.Errorf("get body types scan err : %v", err)
			return nil, 0, err
		}
		bodyTypes = append(bodyTypes, bodyType)
	}

	queryCount := `
			SELECT 
			    COUNT(*) 
			FROM body_types 
			WHERE name ILIKE '%' || $2 || '%'
		`
	errCount := r.client.QueryRow(ctx, queryCount, search).Scan(&count)
	if errCount != nil {
		r.logger.Errorf("get body types count err : %v", err)
		return nil, 0, err
	}
	return bodyTypes, count, nil
}

func (r *BrandPsqlRepository) GetBodyTypeByID(ctx context.Context, id int64) (models.BodyType, error) {
	var bodyType models.BodyType

	query := `
		SELECT
			id, name, image_path
		FROM body_types
		WHERE id = $1
		`

	err := r.client.QueryRow(ctx, query, id).Scan(&bodyType.ID, &bodyType.Name, &bodyType.ImagePath)
	if err != nil {
		r.logger.Errorf("get body type by id query err : %v", err)
		return bodyType, err
	}
	return bodyType, nil
}

func (r *BrandPsqlRepository) UpdateBodyType(ctx context.Context, bodyType models.BodyType) (int64, error) {
	var bodyTypeID int64

	query := `
		UPDATE body_types SET 
		    name = $1, image_path = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id
	`
	err := r.client.QueryRow(ctx, query, bodyType.Name, bodyType.ImagePath, bodyType.ID).Scan(&bodyTypeID)
	if err != nil {
		r.logger.Errorf("update body types err: %v", err)
		return bodyTypeID, err
	}
	return bodyTypeID, nil
}

func (r *BrandPsqlRepository) DeleteBodyType(ctx context.Context, id models.ID) error {
	query := `DELETE FROM body_types WHERE id = $1`
	_, err := r.client.Exec(ctx, query, id.ID)
	if err != nil {
		r.logger.Errorf("delete body types err: %v", err)
		return err
	}
	return nil
}

func (r *BrandPsqlRepository) CreateBrand(ctx context.Context, brand models.Brand) (int64, error) {
	var brandID int64

	tx, err := r.client.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO brands (name, logo_path) VALUES ($1, $2) RETURNING id`

	err = tx.QueryRow(ctx, query, brand.Name, brand.LogoPath).Scan(&brandID)
	if err != nil {
		r.logger.Errorf("create brand err: %v", err)
		return brandID, err
	}

	for _, category := range brand.Categories {
		_, err = tx.Exec(ctx,
			`INSERT INTO brand_category (brand_id, category) VALUES ($1, $2)`,
			brandID, category,
		)
		if err != nil {
			r.logger.Errorf("create brand_categorys err: %v", err)
			return 0, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return brandID, nil
}

func (r *BrandPsqlRepository) GetBrandsByCategory(ctx context.Context, limit, page int64, categoryType, search string) ([]models.Brand, int64, error) {
	var (
		brands []models.Brand
		count  int64
	)

	query := `
		SELECT 
		    b.id, b.name, b.logo_path,
		    ARRAY_AGG(bc.category) AS categories
		FROM brands b
		LEFT JOIN brand_category bc ON bc.brand_id = b.id
		WHERE  bc.category = $1 AND
			b.name ILIKE '%' || $2 || '%'
		GROUP BY b.id
		ORDER BY b.created_at DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.client.Query(ctx, query, categoryType, search, limit, page)
	if err != nil {
		r.logger.Errorf("get brands query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var brand models.Brand
		if err := rows.Scan(&brand.ID, &brand.Name, &brand.LogoPath, &brand.Categories); err != nil {
			r.logger.Errorf("get brands scan err : %v", err)
			return nil, 0, err
		}
		brands = append(brands, brand)
	}

	queryCount := `
			SELECT 
			    COUNT(b.id) 
			FROM brands b
			LEFT JOIN brand_category bc ON bc.brand_id = b.id
			WHERE  bc.category = $1 AND
				b.name ILIKE '%' || $2 || '%'
		`
	errCount := r.client.QueryRow(ctx, queryCount, categoryType, search).Scan(&count)
	if errCount != nil {
		r.logger.Errorf("get brands count err : %v", err)
		return nil, 0, err
	}
	return brands, count, nil
}

func (r *BrandPsqlRepository) UpdateBrand(ctx context.Context, brand models.Brand) (int64, error) {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var id int64

	query := `
		UPDATE brands SET 
		    name = $1, logo_path = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id
	`
	errUpdate := r.client.QueryRow(ctx, query, brand.Name, brand.LogoPath, brand.ID).Scan(&id)
	if errUpdate != nil {
		r.logger.Errorf("update brand err: %v", err)
		return id, errUpdate
	}

	_, err = tx.Exec(ctx, `DELETE FROM brand_category WHERE brand_id = $1`, brand.ID)
	if err != nil {
		r.logger.Errorf("delete old brand_category err: %v", err)
		return id, err
	}

	for _, category := range brand.Categories {
		_, err = tx.Exec(ctx,
			`INSERT INTO brand_category (brand_id, category) VALUES ($1, $2)`,
			brand.ID, category,
		)
		if err != nil {
			r.logger.Errorf("update brand_categorys err: %v", err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
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

func (r *BrandPsqlRepository) DeleteBrandCategory(ctx context.Context, id models.ID) error {
	query := `DELETE FROM brand_category WHERE brand_id = $1 AND category = $2`
	_, err := r.client.Exec(ctx, query, id.ID, id.Category)
	if err != nil {
		r.logger.Errorf("delete brand category err: %v", err)
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
