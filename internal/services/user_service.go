package services

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/repository/storage"
	"context"

	slog "github.com/salamsites/package-log"
)

type UserService struct {
	logger *slog.Logger
	repo   storage.UserRepository
}

func NewUserService(logger *slog.Logger, repo storage.UserRepository) *UserService {
	return &UserService{
		logger: logger,
		repo:   repo,
	}
}

func (s *UserService) GetUsersFromUserService(ctx context.Context, limit, page int64, search string) (dtos.GetUsersResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	users, count, err := s.repo.GetUsersFromUserService(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get all users from user service err: %v", err)
		return dtos.GetUsersResult{}, err
	}
	var dtoUsers []dtos.GetUser
	for _, user := range users {
		dtoUser := dtos.GetUser{
			Id: user.Id,
		}
		if user.FullName.Valid {
			dtoUser.FullName = &user.FullName.String
		}
		if user.Email.Valid {
			dtoUser.Email = &user.Email.String
		}
		if user.PhoneNumber.Valid {
			dtoUser.PhoneNumber = &user.PhoneNumber.String
		}
	}

	result := dtos.GetUsersResult{
		GetUsers: dtoUsers,
		Count:    count,
	}
	return result, nil
}

func (s *UserService) GetUserFirebaseToken(ctx context.Context, userIDs []int64) ([]string, error) {
	tokens, err := s.repo.GetUserFirebaseToken(ctx, userIDs)
	if err != nil {
		s.logger.Errorf("get user firebase token err: %v", err)
		return nil, err
	}
	return tokens, nil
}
