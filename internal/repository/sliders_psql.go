package repository

import (
	"autotm-admin/internal/models"
	"context"
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

	query := `INSERT INTO sliders (image_path, title, platform) VALUES ($1, $2, $3) RETURNING id`

	err := r.client.QueryRow(ctx, query, slider.ImagePath, slider.Title, slider.Platform).Scan(&id)
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
		    id, image_path, title, platform
		FROM sliders
		WHERE platform = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.client.Query(ctx, query, platform, limit, page)
	if err != nil {
		r.logger.Errorf("get sliders query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var slider models.Slider
		if err = rows.Scan(&slider.ID, &slider.ImagePath, &slider.Title, &slider.Platform); err != nil {
			r.logger.Errorf("get sliders scan err : %v", err)
			return nil, 0, err
		}
		sliders = append(sliders, slider)
	}

	queryCount := `
			SELECT 
			    COUNT(*) 
			FROM sliders
			WHERE platform = $1
		`
	errCount := r.client.QueryRow(ctx, queryCount, platform).Scan(&count)
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
		    image_path = $1, title = $2, updated_at = NOW()
		WHERE id = $3 AND platform = $4
		RETURNING id;
	`
	err := r.client.QueryRow(ctx, query, slider.ImagePath, slider.Title, slider.ID, slider.Platform).Scan(&id)
	if err != nil {
		r.logger.Errorf("update slider err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *SliderPsqlRepository) DeleteSlider(ctx context.Context, id models.ID) error {
	query := `DELETE FROM sliders WHERE id = $1`
	_, err := r.client.Exec(ctx, query, id.ID)
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
			id, image_path, platform
		FROM sliders
		WHERE id = $1
	`
	err := r.client.QueryRow(ctx, query, id).Scan(&slider.ID, &slider.ImagePath, &slider.Platform)
	if err != nil {
		r.logger.Errorf("get slider by id query err : %v", err)
		return slider, err
	}
	return slider, nil
}
