package repository

import (
	"autotm-admin/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	slog "github.com/salamsites/package-log"
	spsql "github.com/salamsites/package-psql"
)

type SettingsPsqlRepository struct {
	logger *slog.Logger
	client spsql.Client
}

func NewSettingsPsqlRepository(logger *slog.Logger, client spsql.Client) *SettingsPsqlRepository {
	return &SettingsPsqlRepository{
		logger: logger,
		client: client,
	}
}

func (r *SettingsPsqlRepository) CreateRole(ctx context.Context, role models.Role) (int64, error) {
	var id int64

	query := `INSERT INTO roles (name, role) VALUES (@name, @role) RETURNING id`

	args := pgx.NamedArgs{
		"name": role.Name,
		"role": role.Role,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("create role err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *SettingsPsqlRepository) GetRoleByID(ctx context.Context, roleID int64) (models.Role, error) {
	var role models.Role

	query := `
		SELECT
            id, name, role
        FROM roles
		WHERE id = @id
    `

	args := pgx.NamedArgs{
		"id": roleID,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&role.ID, &role.Name, &role.Role)
	if err != nil {
		r.logger.Errorf("get role err: %v", err)
		return role, err
	}

	return role, nil
}

func (r *SettingsPsqlRepository) GetAllRoles(ctx context.Context, limit, page int64, search string) ([]models.Role, int64, error) {
	var (
		roles []models.Role
		count int64
	)

	query := `
		SELECT 
		    id, name, role 
		FROM roles
		WHERE name ILIKE '%' || @search || '%'
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
		r.logger.Errorf("get all roles query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Role); err != nil {
			r.logger.Errorf("get all roles scan err : %v", err)
			return nil, 0, err
		}
		roles = append(roles, role)
	}

	queryCount := `
			SELECT 
			    COUNT(*) 
			FROM roles
			WHERE name ILIKE '%' || @search || '%'
		`

	argsCount := pgx.NamedArgs{
		"search": search,
	}
	errCount := r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
	if errCount != nil {
		r.logger.Errorf("get all roles count err : %v", err)
		return nil, 0, err
	}
	return roles, count, nil
}

func (r *SettingsPsqlRepository) UpdateRole(ctx context.Context, role models.Role) (int64, error) {
	var id int64

	query := `
		UPDATE roles SET 
		    name = @name, role = @role, updated_at = NOW()
		WHERE id = @id
		RETURNING id
	`

	args := pgx.NamedArgs{
		"name": role.Name,
		"role": role.Role,
		"id":   role.ID,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("update role err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *SettingsPsqlRepository) DeleteRole(ctx context.Context, id int64) error {
	query := `DELETE FROM roles WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}
	_, err := r.client.Exec(ctx, query, args)
	if err != nil {
		r.logger.Errorf("delete role err: %v", err)
		return err
	}
	return nil
}

// Admin Users
func (r *SettingsPsqlRepository) CreateUser(ctx context.Context, user models.User) (int64, error) {
	var id int64

	query := `INSERT INTO admin_users (username, login, password, role_id) VALUES (@username, @login, @password, @role_id) RETURNING id`

	args := pgx.NamedArgs{
		"username": user.Username,
		"login":    user.Login,
		"password": user.Password,
		"role_id":  user.RoleID,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("create user err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *SettingsPsqlRepository) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	var (
		user models.User
	)

	query := `
		SELECT
			id, username, login, password, role_id
		FROM admin_users
		WHERE login = @login
		LIMIT 1
	`

	args := pgx.NamedArgs{
		"login": login,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&user.ID, &user.Username, &user.Login, &user.Password, &user.RoleID)
	if err != nil {
		r.logger.Errorf("get user by login err: %v", err)
		return user, err
	}
	return user, nil
}

func (r *SettingsPsqlRepository) GetAllUsers(ctx context.Context, limit, page int64, search string) ([]models.User, int64, error) {
	var (
		users []models.User
		count int64
	)

	query := `
		SELECT 
		    u.id, u.username, u.login, u.password, 
		    u.role_id, r.name AS role_name 
		FROM admin_users u
			LEFT JOIN roles r ON u.role_id = r.id
		WHERE (u.username ILIKE '%' || @search || '%' OR r.name ILIKE '%' || @search || '%' OR u.login ILIKE '%' || @search || '%')
		ORDER BY u.created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	args := pgx.NamedArgs{
		"search": search,
		"limit":  limit,
		"offset": page,
	}

	rows, err := r.client.Query(ctx, query, args)
	if err != nil {
		r.logger.Errorf("get all users query err : %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Login, &user.Password, &user.RoleID, &user.RoleName); err != nil {
			r.logger.Errorf("get all users scan err : %v", err)
			return nil, 0, err
		}
		users = append(users, user)
	}

	queryCount := `
			SELECT 
			    COUNT(u.id) 
			FROM admin_users u
				LEFT JOIN roles r ON u.role_id = r.id
			WHERE (u.username ILIKE '%' || @search || '%' OR r.name ILIKE '%' || @search || '%' OR u.login ILIKE '%' || @search || '%')
		`

	argsCount := pgx.NamedArgs{
		"search": search,
	}
	errCount := r.client.QueryRow(ctx, queryCount, argsCount).Scan(&count)
	if errCount != nil {
		r.logger.Errorf("get all users count err : %v", err)
		return nil, 0, err
	}
	return users, count, nil
}

func (r *SettingsPsqlRepository) UpdateUser(ctx context.Context, user models.User) (int64, error) {
	var id int64

	query := `
		UPDATE admin_users SET 
		    username = @username, login = @login, password = @password, role_id = @role_id, updated_at = NOW()
		WHERE id = @id
		RETURNING id
	`
	args := pgx.NamedArgs{
		"username": user.Username,
		"login":    user.Login,
		"password": user.Password,
		"role_id":  user.RoleID,
		"id":       user.ID,
	}
	err := r.client.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		r.logger.Errorf("update user err: %v", err)
		return id, err
	}
	return id, nil
}

func (r *SettingsPsqlRepository) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM admin_users WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}
	_, err := r.client.Exec(ctx, query, args)
	if err != nil {
		r.logger.Errorf("delete users err: %v", err)
		return err
	}
	return nil
}
