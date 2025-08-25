package repository

import (
	"autotm-admin/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
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
		WHERE (full_name ILIKE '%' || @search || '%' OR mail ILIKE '%' || @search || '%' OR phone_number ILIKE '%' || @search || '%')
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
			WHERE (full_name ILIKE '%' || @search || '%' OR mail ILIKE '%' || @search || '%' OR phone_number ILIKE '%' || @search || '%')
		`

	argsCount := pgx.NamedArgs{
		"search": search,
	}
	errCount := r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
	if errCount != nil {
		r.logger.Errorf("get all users from user service count err : %v", err)
		return nil, 0, err
	}
	return users, count, nil
}

func (r *UserPsqlRepository) GetUserFirebaseToken(ctx context.Context, userId int64) (string, error) {
	var (
		token string
	)

	query := `
 		SELECT
			firebase_token
        FROM user_devices
 		WHERE user_id = @user_id;
	`

	args := pgx.NamedArgs{
		"user_id": userId,
	}

	err := r.client.QueryRow(ctx, query, args).Scan(&token)
	if err != nil {
		return "", err
	}

	return token, nil
}
