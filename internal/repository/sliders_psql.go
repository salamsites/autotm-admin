package repository

import (
	"autotm-admin/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	slog "github.com/salamsites/package-log"
	spsql "github.com/salamsites/package-psql"
)

type SliderPsqlRepository struct {
	logger *slog.Logger
	client spsql.Client
}

func NewSliderPsqlRepository(logger *slog.Logger, client spsql.Client) *SliderPsqlRepository {
	return &SliderPsqlRepository{
		logger: logger,
		client: client,
	}
}

func (r *SliderPsqlRepository) CreateSlider(ctx context.Context, slider models.Slider) (int64, error) {
	var id int64

	query := `
			INSERT INTO sliders 
    			(image_path_tm, image_path_en, image_path_ru, upload_id_tm, upload_id_en, upload_id_ru, platform) 
			VALUES (@image_path_tm, @image_path_en, @image_path_ru, @upload_id_tm, @upload_id_en, @upload_id_ru, @platform) 
			RETURNING id
		`

	args := pgx.NamedArgs{
		"image_path_tm": slider.ImagePathTM,
		"image_path_en": slider.ImagePathEN,
		"image_path_ru": slider.ImagePathRU,
		"upload_id_tm":  slider.UploadIdTM,
		"upload_id_en":  slider.UploadIdEN,
		"upload_id_ru":  slider.UploadIdRU,
		"platform":      slider.Platform,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("create err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *SliderPsqlRepository) GetAllSliders(ctx context.Context, limit, page int64, platform string) ([]models.Slider, int64, error) {
	var (
		sliders []models.Slider
		count   int64
	)

	query := `
		SELECT 
		    id, image_path_tm, image_path_en, image_path_ru, 
		    upload_id_tm, upload_id_en, upload_id_ru, platform
		FROM sliders
		WHERE platform = @platform
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	args := pgx.NamedArgs{
		"platform": platform,
		"limit":    limit,
		"offset":   page,
	}
	rows, err := r.client.Query(ctx, query, args)
	if err != nil {
		r.logger.Errorf("get sliders query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var slider models.Slider
		if err = rows.Scan(&slider.ID, &slider.ImagePathTM, &slider.ImagePathEN, &slider.ImagePathRU,
			&slider.UploadIdTM, &slider.UploadIdEN, &slider.UploadIdRU, &slider.Platform,
		); err != nil {
			r.logger.Errorf("get sliders scan err : %v", err)
			return nil, 0, err
		}
		sliders = append(sliders, slider)
	}

	queryCount := `
			SELECT 
			    COUNT(*) 
			FROM sliders
			WHERE platform = @platform
		`

	argsCount := pgx.NamedArgs{
		"platform": platform,
	}
	errCount := r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
	if errCount != nil {
		r.logger.Errorf("get sliders count err : %v", err)
		return nil, 0, err
	}
	return sliders, count, nil
}

func (r *SliderPsqlRepository) UpdateSlider(ctx context.Context, slider models.Slider) (int64, error) {
	var id int64

	query := `
		UPDATE sliders SET 
		    image_path_tm = @image_path_tm, image_path_en = @image_path_en, image_path_ru = @image_path_ru, platform = @platform, 
		    upload_id_tm = @upload_id_tm, upload_id_en = @upload_id_en, upload_id_ru = @upload_id_ru, updated_at = NOW()
		WHERE id = @id
		RETURNING id;
	`

	args := pgx.NamedArgs{
		"image_path_tm": slider.ImagePathTM,
		"image_path_en": slider.ImagePathEN,
		"image_path_ru": slider.ImagePathRU,
		"platform":      slider.Platform,
		"upload_id_tm":  slider.UploadIdTM,
		"upload_id_en":  slider.UploadIdEN,
		"upload_id_ru":  slider.UploadIdRU,
		"id":            slider.ID,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("update slider err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *SliderPsqlRepository) DeleteSlider(ctx context.Context, id int64) error {
	query := `DELETE FROM sliders WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}
	_, err := r.client.Exec(ctx, query, args)
	if err != nil {
		r.logger.Errorf("delete slider err: %v", err)
		return err
	}
	return nil
}

func (r *SliderPsqlRepository) GetSliderByID(ctx context.Context, id int64) (models.Slider, error) {
	var slider models.Slider

	query := `
		SELECT
			id, image_path_tm, image_path_en, image_path_ru, 
			upload_id_tm, upload_id_en, upload_id_ru, platform
		FROM sliders
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&slider.ID, &slider.ImagePathTM, &slider.ImagePathEN,
		&slider.ImagePathRU, &slider.UploadIdTM, &slider.UploadIdEN, &slider.UploadIdRU, &slider.Platform,
	)
	if err != nil {
		r.logger.Errorf("get slider by id query err : %v", err)
		return slider, err
	}
	return slider, nil
}
