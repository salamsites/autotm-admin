package storage

import (
	"autotm-admin/internal/models"
	"context"
)

type UserRepository interface {
	GetUsersFromUserService(ctx context.Context, limit, page int64, search string) ([]models.GetUser, int64, error)
	GetUserFirebaseToken(ctx context.Context, userId int64) (string, error)
}
