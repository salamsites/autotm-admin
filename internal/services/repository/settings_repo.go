package repository

import (
	"autotm-admin/internal/dtos"
	"context"
)

type SettingsService interface {
	// Role
	CreateRole(ctx context.Context, role dtos.CreateRoleReq) (int64, error)
	GetRoleByID(ctx context.Context, roleID int64) (dtos.Role, error)
	GetAllRoles(ctx context.Context, limit, page int64, search string) (dtos.RoleResult, error)
	UpdateRole(ctx context.Context, role dtos.UpdateRoleReq) (int64, error)
	DeleteRole(ctx context.Context, id int64) error

	// User
	CreateUser(ctx context.Context, user dtos.CreateUserReq) (int64, error)
	InitSuperAdmin(ctx context.Context) error
	GetAllUsers(ctx context.Context, limit, page int64, search string) (dtos.UserResult, error)
	UpdateUser(ctx context.Context, user dtos.UpdateUserReq) (int64, error)
	DeleteUser(ctx context.Context, id int64) error
}
