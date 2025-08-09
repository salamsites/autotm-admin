package repository

import (
	"autotm-admin/internal/models"
	"context"
	slog "github.com/salamsites/package-log"
	spsql "github.com/salamsites/package-psql"
)

type UserPsqlRepository struct {
	logger *slog.Logger
	client spsql.Client
}

func NewUserPsqlRepository(logger *slog.Logger, client spsql.Client) *UserPsqlRepository {
	return &UserPsqlRepository{
		logger: logger,
		client: client,
	}
}

func (r *UserPsqlRepository) GetUsersFromUserService(ctx context.Context, limit, page int64, search string) ([]models.GetUser, int64, error) {
	var (
		users []models.GetUser
		count int64
	)

	query := `
		SELECT 
		    id, full_name, mail, phone_number
		FROM users 
		WHERE (full_name ILIKE '%' || $1 || '%' OR mail ILIKE '%' || $1 || '%' OR phone_number ILIKE '%' || $1 || '%')
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.client.Query(ctx, query, search, limit, page)
	if err != nil {
		r.logger.Errorf("get all users from user service query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.GetUser
		if err = rows.Scan(&user.Id, &user.FullName, &user.Email, &user.PhoneNumber); err != nil {
			r.logger.Errorf("get all users from user service scan err : %v", err)
			return nil, 0, err
		}
		users = append(users, user)
	}

	queryCount := `
			SELECT 
			    COUNT(*) 
			FROM users
			WHERE (full_name ILIKE '%' || $1 || '%' OR mail ILIKE '%' || $1 || '%' OR phone_number ILIKE '%' || $1 || '%')
		`
	errCount := r.client.QueryRow(ctx, queryCount, search).Scan(&count)
	if errCount != nil {
		r.logger.Errorf("get all users from user service count err : %v", err)
		return nil, 0, err
	}
	return users, count, nil
}
