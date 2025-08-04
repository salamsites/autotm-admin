package storage

import (
	"autotm-admin/internal/models"
	"context"
)

type SettingsRepository interface {
	CreateRole(ctx context.Context, model models.Role) (int64, error)
	GetRoleByID(ctx context.Context, roleID int64) (models.Role, error)
	GetAllRoles(ctx context.Context, limit, page int64, search string) ([]models.Role, int64, error)
	UpdateRole(ctx context.Context, role models.Role) (int64, error)
	DeleteRole(ctx context.Context, id models.ID) error

	// Users
	CreateUser(ctx context.Context, model models.User) (int64, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	GetAllUsers(ctx context.Context, limit, page int64, search string) ([]models.User, int64, error)
	UpdateUser(ctx context.Context, role models.User) (int64, error)
	DeleteUser(ctx context.Context, id models.ID) error
}
