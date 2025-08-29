package repository

import (
	"autotm-admin/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
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

	query := ` 
				INSERT INTO body_types 
				    (name_tm, name_en, name_ru, image_path, category, upload_id) 
 				VALUES 
 				    (@name_tm, @name_en, @name_ru, @image_path, @category, @upload_id) 
 				RETURNING id 
		`

	args := pgx.NamedArgs{
		"name_tm":    bodyType.NameTM,
		"name_en":    bodyType.NameEN,
		"name_ru":    bodyType.NameRU,
		"image_path": bodyType.ImagePath,
		"category":   bodyType.Category,
		"upload_id":  bodyType.UploadId,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&id)
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
				id, name_tm, name_en, name_ru, category, image_path, upload_id
            FROM body_types
			WHERE category = @category AND 
			    (name_tm ILIKE '%' || @search || '%' OR name_en ILIKE '%' || @search || '%' OR name_ru ILIKE '%' || @search || '%')
			ORDER BY created_at DESC
			LIMIT @limit OFFSET @offset;
		`

	args := pgx.NamedArgs{
		"category": category,
		"search":   search,
		"limit":    limit,
		"offset":   page,
	}

	rows, err := r.client.Query(ctx, query, args)
	if err != nil {
		r.logger.Errorf("get body types query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var bodyType models.BodyType
		if err = rows.Scan(&bodyType.ID, &bodyType.NameTM, &bodyType.NameEN, &bodyType.NameRU, &bodyType.Category, &bodyType.ImagePath, &bodyType.UploadId); err != nil {
			r.logger.Errorf("get body types scan err : %v", err)
			return nil, 0, err
		}
		bodyTypes = append(bodyTypes, bodyType)
	}

	queryCount := `
			SELECT 
			    COUNT(*) 
			FROM body_types 
			WHERE category = @category AND 
				(name_tm ILIKE '%' || @search || '%' OR name_en ILIKE '%' || @search || '%' OR name_ru ILIKE '%' || @search || '%')	
		`

	argsCount := pgx.NamedArgs{
		"category": category,
		"search":   search,
	}
	errCount := r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
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
			id, name_tm, name_en, name_ru, 
			image_path, category, upload_id
		FROM body_types
		WHERE id = @id
		`

	args := pgx.NamedArgs{
		"id": id,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&bodyType.ID, &bodyType.NameTM, &bodyType.NameEN, &bodyType.NameRU, &bodyType.ImagePath, &bodyType.Category, &bodyType.ImagePath, &bodyType.UploadId)
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
		    name_tm = @name_tm, name_en = @name_en, name_ru = @name_ru, image_path = @image_path, 
		    category = @category, upload_id = @upload_id, updated_at = NOW()
		WHERE id = @id
		RETURNING id
	`

	args := pgx.NamedArgs{
		"name_tm":    bodyType.NameTM,
		"name_en":    bodyType.NameEN,
		"name_ru":    bodyType.NameRU,
		"image_path": bodyType.ImagePath,
		"category":   bodyType.Category,
		"upload_id":  bodyType.UploadId,
		"id":         bodyType.ID,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&bodyTypeID)
	if err != nil {
		r.logger.Errorf("update body types err: %v", err)
		return bodyTypeID, err
	}
	return bodyTypeID, nil
}

func (r *BrandPsqlRepository) DeleteBodyType(ctx context.Context, id int64) error {
	query := `DELETE FROM body_types WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}
	_, err := r.client.Exec(ctx, query, args)
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

	query := ` 
		INSERT INTO brands 
     				(name, logo_path, upload_id) 
 		VALUES (@name, @logo_path, @upload_id) 
 		RETURNING id 
	`

	args := pgx.NamedArgs{
		"name":      brand.Name,
		"logo_path": brand.LogoPath,
		"upload_id": brand.UploadId,
	}

	err = tx.QueryRow(ctx, query, args).Scan(&brandID)
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

func (r *BrandPsqlRepository) GetBrands(ctx context.Context, limit, page int64, category, search string) ([]models.Brand, int64, error) {
	var (
		brands []models.Brand
		count  int64
	)

	query := `
		SELECT 
		    b.id, b.name, b.logo_path, b.upload_id,
		    ARRAY_AGG(bc.category) AS categories
		FROM brands b
		LEFT JOIN brand_categories bc ON bc.brand_id = b.id
		WHERE  bc.category = @category AND
			b.name ILIKE '%' || @search || '%'
		GROUP BY b.id
		ORDER BY b.created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	args := pgx.NamedArgs{
		"category": category,
		"search":   search,
		"limit":    limit,
		"offset":   page,
	}

	rows, err := r.client.Query(ctx, query, args)
	if err != nil {
		r.logger.Errorf("get brands query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var brand models.Brand
		if err := rows.Scan(&brand.ID, &brand.Name, &brand.LogoPath, &brand.UploadId, &brand.Categories); err != nil {
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
			WHERE  bc.category = @category AND
				b.name ILIKE '%' || @search || '%'
		`

	argsCount := pgx.NamedArgs{
		"category": category,
		"search":   search,
	}
	errCount := r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
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
		    name = @name, logo_path = @logo_path, upload_id = @upload_id, updated_at = NOW()
		WHERE id = @id
		RETURNING id
	`

	args := pgx.NamedArgs{
		"name":      brand.Name,
		"logo_path": brand.LogoPath,
		"upload_id": brand.UploadId,
		"id":        brand.ID,
	}
	errUpdate := tx.QueryRow(ctx, query, args).Scan(&id)
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
			id, name, logo_path, upload_id
		FROM brands
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&brand.ID, &brand.Name, &brand.LogoPath, &brand.UploadId)
	if err != nil {
		r.logger.Errorf("get brand by id query err : %v", err)
		return brand, err
	}
	return brand, nil
}

func (r *BrandPsqlRepository) DeleteBrandCategory(ctx context.Context, id int64, category string) error {
	query := `DELETE FROM brand_categories WHERE brand_id = @brand_id AND category = @category`

	args := pgx.NamedArgs{
		"brand_id": id,
		"category": category,
	}
	_, err := r.client.Exec(ctx, query, args)
	if err != nil {
		r.logger.Errorf("delete brand category err: %v", err)
		return err
	}
	return nil
}

func (r *BrandPsqlRepository) CreateModel(ctx context.Context, model models.Model) (int64, error) {
	var id int64

	query := ` INSERT INTO models (name, brand_id, category) VALUES (@name, @brand_id, @category) RETURNING id `

	args := pgx.NamedArgs{
		"name":     model.Name,
		"brand_id": model.BrandID,
		"category": model.Category,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&id)
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
		    m.id, m.name, b.logo_path, b.upload_id, 
		    m.brand_id, b.name, m.category
		FROM models m
			LEFT JOIN brands b ON m.brand_id = b.id
		WHERE  m.category = @category AND 
		    ( m.name ILIKE '%' || @search || '%' OR b.name ILIKE '%' || @search || '%' )
		ORDER BY m.created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	args := pgx.NamedArgs{
		"category": category,
		"search":   search,
		"limit":    limit,
		"offset":   page,
	}

	rows, err := r.client.Query(ctx, query, args)
	if err != nil {
		r.logger.Errorf("get models query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var brandModel models.Model
		if err = rows.Scan(&brandModel.ID, &brandModel.Name, &brandModel.LogoPath,
			&brandModel.UploadId, &brandModel.BrandID, &brandModel.BrandName, &brandModel.Category,
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
		WHERE  m.category = @category AND 
		    ( m.name ILIKE '%' || @search || '%' OR b.name ILIKE '%' || @search || '%' )
		`

	argsCount := pgx.NamedArgs{
		"category": category,
		"search":   search,
	}
	errCount := r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
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
		    name = @name, brand_id = @brand_id, category = @category, updated_at = NOW()
		WHERE id = @id
		RETURNING id;
	`

	args := pgx.NamedArgs{
		"name":     model.Name,
		"brand_id": model.BrandID,
		"category": model.Category,
		"id":       model.ID,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("update brand model err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *BrandPsqlRepository) DeleteModel(ctx context.Context, id int64) error {
	query := `DELETE FROM models WHERE id = @id`

	args := pgx.NamedArgs{
		"id": id,
	}
	_, err := r.client.Exec(ctx, query, args)
	if err != nil {
		r.logger.Errorf("delete model err: %v", err)
		return err
	}
	return nil
}

func (r *BrandPsqlRepository) CreateDescription(ctx context.Context, description models.Description) (int64, error) {
	var id int64

	query := ` INSERT INTO descriptions (name_tm, name_en, name_ru, category) VALUES (@name_tm, @name_en, @name_ru, @category) RETURNING id `

	args := pgx.NamedArgs{
		"name_tm": description.NameTM,
		"name_en": description.NameEN,
		"name_ru": description.NameRU,
		"type":    description.Category,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("Error creating description: %s", err.Error())
		return id, err
	}
	return id, nil
}

func (r *BrandPsqlRepository) GetDescriptions(ctx context.Context, limit, page int64, search, category string) ([]models.Description, int64, error) {
	var (
		descriptions []models.Description
		count        int64
	)

	query := `
			SELECT 
				id, name_tm, name_en, name_ru, category
            FROM descriptions
			WHERE category = @category
			    (name_tm ILIKE '%' || @search || '%' OR name_en ILIKE '%' || @search || '%' OR name_ru ILIKE '%' || @search || '%')
			ORDER BY created_at DESC
			LIMIT @limit OFFSET @offset;
		`

	args := pgx.NamedArgs{
		"type":   category,
		"search": search,
		"limit":  limit,
		"offset": page,
	}
	rows, err := r.client.Query(ctx, query, args)
	if err != nil {
		r.logger.Errorf("get descriptions query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var description models.Description
		if err = rows.Scan(&description.ID, &description.NameTM, &description.NameEN, &description.NameRU, &description.Category); err != nil {
			r.logger.Errorf("get descriptions scan err : %v", err)
			return nil, 0, err
		}
		descriptions = append(descriptions, description)
	}

	queryCount := `
			SELECT 
			    COUNT(*) 
			FROM descriptions 
			WHERE category = @category
				(name_tm ILIKE '%' || @search || '%' OR name_en ILIKE '%' || @search || '%' OR name_ru ILIKE '%' || @search || '%')	
		`

	argsCount := pgx.NamedArgs{
		"type":   category,
		"search": search,
	}
	errCount := r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
	if errCount != nil {
		r.logger.Errorf("get descriptions count err : %v", err)
		return nil, 0, err
	}
	return descriptions, count, nil
}

func (r *BrandPsqlRepository) UpdateDescription(ctx context.Context, description models.Description) (int64, error) {
	var descriptionID int64

	query := `
		UPDATE descriptions SET 
		    name_tm = @name_tm, name_en = @name_en, name_ru = @name_ru, 
		    category = @type, updated_at = NOW()
		WHERE id = @id
		RETURNING id
	`

	args := pgx.NamedArgs{
		"name_tm": description.NameTM,
		"name_en": description.NameEN,
		"name_ru": description.NameRU,
		"type":    description.Category,
		"id":      description.ID,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&descriptionID)
	if err != nil {
		r.logger.Errorf("update descriptions err: %v", err)
		return descriptionID, err
	}
	return descriptionID, nil
}

func (r *BrandPsqlRepository) DeleteDescription(ctx context.Context, id int64) error {
	query := `DELETE FROM descriptions WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}
	_, err := r.client.Exec(ctx, query, args)
	if err != nil {
		r.logger.Errorf("delete descriptions err: %v", err)
		return err
	}
	return nil
}
