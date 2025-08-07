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

	query := ` INSERT INTO body_types (name_tm, name_en, name_ru, image_path, category) VALUES ($1, $2, $3, $4, $5) RETURNING id `

	err := r.client.QueryRow(ctx, query, bodyType.NameTM, bodyType.NameEN, bodyType.NameRU, bodyType.ImagePath, bodyType.Category).Scan(&id)
	if err != nil {
		r.logger.Errorf("Error creating body type: %s", err.Error())
		return id, err
	}
	return id, nil
}

func (r *BrandPsqlRepository) GetBodyType(ctx context.Context, limit, page int64, category, search string) ([]models.BodyType, int64, error) {
	var (
		bodyTypes []models.BodyType
		count     int64
	)

	query := `
			SELECT 
				id, name_tm, name_en, name_ru, category, image_path
            FROM body_types
			WHERE category = $1 AND 
			    (name_tm ILIKE '%' || $2 || '%' OR name_en ILIKE '%' || $2 || '%' OR name_ru ILIKE '%' || $2 || '%')
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4;
		`

	rows, err := r.client.Query(ctx, query, category, search, limit, page)
	if err != nil {
		r.logger.Errorf("get body types query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var bodyType models.BodyType
		if err = rows.Scan(&bodyType.ID, &bodyType.NameTM, &bodyType.NameEN, &bodyType.NameRU, &bodyType.Category, &bodyType.ImagePath); err != nil {
			r.logger.Errorf("get body types scan err : %v", err)
			return nil, 0, err
		}
		bodyTypes = append(bodyTypes, bodyType)
	}

	queryCount := `
			SELECT 
			    COUNT(*) 
			FROM body_types 
			WHERE category = $1 AND 
				(name_tm ILIKE '%' || $2 || '%' OR name_en ILIKE '%' || $2 || '%' OR name_ru ILIKE '%' || $2 || '%')	
		`
	errCount := r.client.QueryRow(ctx, queryCount, category, search).Scan(&count)
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
			id, name_tm, name_en, name_ru, image_path, category
		FROM body_types
		WHERE id = $1
		`

	err := r.client.QueryRow(ctx, query, id).Scan(&bodyType.ID, &bodyType.NameTM, &bodyType.NameEN, &bodyType.NameRU, &bodyType.ImagePath, &bodyType.Category)
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
		    name_tm = $1, name_en = $2, name_ru = $3, image_path = $4, category = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING id
	`
	err := r.client.QueryRow(ctx, query, bodyType.NameTM, bodyType.NameEN, bodyType.NameRU, bodyType.ImagePath, bodyType.Category, bodyType.ID).Scan(&bodyTypeID)
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
			`INSERT INTO brand_categories (brand_id, category) VALUES ($1, $2)`,
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

func (r *BrandPsqlRepository) GetBrands(ctx context.Context, limit, page int64, categoryType, search string) ([]models.Brand, int64, error) {
	var (
		brands []models.Brand
		count  int64
	)

	query := `
		SELECT 
		    b.id, b.name, b.logo_path,
		    ARRAY_AGG(bc.category) AS categories
		FROM brands b
		LEFT JOIN brand_categories bc ON bc.brand_id = b.id
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
			LEFT JOIN brand_categories bc ON bc.brand_id = b.id
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
		return 0, errUpdate
	}

	_, err = tx.Exec(ctx, `DELETE FROM brand_categories WHERE brand_id = $1`, brand.ID)
	if err != nil {
		r.logger.Errorf("delete old brand_category err: %v", err)
		return 0, err
	}

	for _, category := range brand.Categories {
		_, err = tx.Exec(ctx,
			`INSERT INTO brand_categories (brand_id, category) VALUES ($1, $2)`,
			brand.ID, category,
		)
		if err != nil {
			r.logger.Errorf("update brand_categorys err: %v", err)
			return 0, err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
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
	query := `DELETE FROM brand_categories WHERE brand_id = $1 AND category = $2`
	_, err := r.client.Exec(ctx, query, id.ID, id.Category)
	if err != nil {
		r.logger.Errorf("delete brand category err: %v", err)
		return err
	}
	return nil
}

func (r *BrandPsqlRepository) CreateModel(ctx context.Context, model models.Model) (int64, error) {
	var id int64

	query := ` INSERT INTO models (name, brand_id, category) VALUES ($1, $2, $3) RETURNING id `

	err := r.client.QueryRow(ctx, query, model.Name, model.BrandID, model.Category).Scan(&id)
	if err != nil {
		r.logger.Errorf("create model err : %v", err)
		return id, err
	}
	return id, nil
}

func (r *BrandPsqlRepository) GetModels(ctx context.Context, limit, page int64, category, search string) ([]models.Model, int64, error) {
	var (
		brandModels []models.Model
		count       int64
	)

	query := `
		SELECT 
		    m.id, m.name, b.logo_path, m.brand_id, b.name,
		    m.category
		FROM models m
			LEFT JOIN brands b ON m.brand_id = b.id
		WHERE  m.category = $1 AND 
		    ( m.name ILIKE '%' || $2 || '%' OR b.name ILIKE '%' || $2 || '%' )
		ORDER BY m.created_at DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.client.Query(ctx, query, category, search, limit, page)
	if err != nil {
		r.logger.Errorf("get models query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var brandModel models.Model
		if err = rows.Scan(&brandModel.ID, &brandModel.Name, &brandModel.LogoPath,
			&brandModel.BrandID, &brandModel.BrandName, &brandModel.Category,
		); err != nil {
			r.logger.Errorf("get models scan err : %v", err)
			return nil, 0, err
		}
		brandModels = append(brandModels, brandModel)
	}

	queryCount := `
		SELECT 
			COUNT(m.id) 
		FROM models m
			LEFT JOIN brands b ON m.brand_id = b.id
		WHERE  m.category = $1 AND 
		    ( m.name ILIKE '%' || $2 || '%' OR b.name ILIKE '%' || $2 || '%' )
		`
	errCount := r.client.QueryRow(ctx, queryCount, category, search).Scan(&count)
	if errCount != nil {
		r.logger.Errorf("get models count err : %v", err)
		return nil, 0, err
	}
	return brandModels, count, nil
}

func (r *BrandPsqlRepository) UpdateModel(ctx context.Context, model models.Model) (int64, error) {
	var id int64

	query := `
		UPDATE models SET 
		    name = $1, brand_id = $2, category = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id;
	`
	err := r.client.QueryRow(ctx, query, model.Name, model.BrandID, model.Category, model.ID).Scan(&id)
	if err != nil {
		r.logger.Errorf("update brand model err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *BrandPsqlRepository) DeleteModel(ctx context.Context, id models.ID) error {
	query := `DELETE FROM models WHERE id = $1`
	_, err := r.client.Exec(ctx, query, id.ID)
	if err != nil {
		r.logger.Errorf("delete model err: %v", err)
		return err
	}
	return nil
}
