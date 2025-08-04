package services

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"autotm-admin/internal/models"
	"autotm-admin/internal/repository/storage"
	"autotm-admin/internal/services/repository"
	"context"
	slog "github.com/salamsites/package-log"
)

type AutoStoreService struct {
	logger      *slog.Logger
	repo        storage.AutoStoreRepository
	userService repository.UserService
}

func NewAutoStoreService(logger *slog.Logger, repo storage.AutoStoreRepository, userService repository.UserService) *AutoStoreService {
	return &AutoStoreService{
		logger:      logger,
		repo:        repo,
		userService: userService,
	}
}

func (s *AutoStoreService) CreateAutoStore(ctx context.Context, autoStore dtos.CreateAutoStoreReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(autoStore); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	newAutoStore := models.AutoStore{
		UserID:      autoStore.UserID,
		PhoneNumber: autoStore.PhoneNumber,
		Email:       autoStore.Email,
		StoreName:   autoStore.StoreName,
		Images:      autoStore.Images,
		LogoPath:    autoStore.LogoPath,
		RegionID:    autoStore.RegionID,
		CityID:      autoStore.CityID,
	}

	autoStoreID, err := s.repo.CreateAutoStore(ctx, newAutoStore)
	if err != nil {
		s.logger.Errorf("create err: %v", err)
		return autoStoreID, err
	}
	return autoStoreID, nil
}

func (s *AutoStoreService) GetUsersFromUserService(ctx context.Context, limit, page int64, search string) (dtos.GetUserResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	users, count, err := s.userService.GetUsers(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get users err: %v", err)
		return dtos.GetUserResult{}, err
	}

	result := dtos.GetUserResult{
		Users: users,
		Count: count,
	}

	return result, nil
}

func (s *AutoStoreService) GetAutoStores(ctx context.Context, limit, page int64, search string) (dtos.AutoStoresResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	autoStores, count, err := s.repo.GetAutoStores(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get autoStores err: %v", err)
		return dtos.AutoStoresResult{}, err
	}

	userIDMap := make(map[int64]struct{})
	var id dtos.GetUserByIDsReq
	for _, a := range autoStores {
		if _, exists := userIDMap[a.UserID]; !exists {
			userIDMap[a.UserID] = struct{}{}
			id.Ids = append(id.Ids, a.UserID)
		}
	}

	users, err := s.userService.GetUserByIds(ctx, id)
	if err != nil {
		s.logger.Errorf("get user by ids err: %v", err)
		return dtos.AutoStoresResult{}, err
	}

	userMap := make(map[int64]dtos.GetUsers)
	for _, user := range users {
		userMap[user.Id] = user
	}

	var dtoAutoStores []dtos.AutoStore
	for _, autoStore := range autoStores {
		user := userMap[autoStore.UserID]
		dtoAutoStores = append(dtoAutoStores, dtos.AutoStore{
			ID:          autoStore.ID,
			PhoneNumber: autoStore.PhoneNumber,
			Email:       autoStore.Email,
			StoreName:   autoStore.StoreName,
			Images:      autoStore.Images,
			LogoPath:    autoStore.LogoPath,
			RegionID:    autoStore.RegionID,
			CityID:      autoStore.CityID,
			UserID:      autoStore.UserID,
			UserName:    user.FullName,
		})
	}

	result := dtos.AutoStoresResult{
		AutoStores: dtoAutoStores,
		Count:      count,
	}

	return result, nil
}

func (s *AutoStoreService) UpdateAutoStore(ctx context.Context, autoStore dtos.UpdateAutoStoreReq) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(autoStore); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	newAutoStore := models.AutoStore{
		ID:          autoStore.ID,
		PhoneNumber: autoStore.PhoneNumber,
		Email:       autoStore.Email,
		StoreName:   autoStore.StoreName,
		Images:      autoStore.Images,
		LogoPath:    autoStore.LogoPath,
		RegionID:    autoStore.RegionID,
		CityID:      autoStore.CityID,
		Address:     autoStore.Address,
	}

	autoStoreID, err := s.repo.UpdateAutoStore(ctx, newAutoStore)
	if err != nil {
		s.logger.Errorf("update auto store err: %v", err)
		return id, err
	}
	id.ID = autoStoreID
	return id, nil
}

func (s *AutoStoreService) DeleteAutoStore(ctx context.Context, id int64) error {
	deleteID := models.ID{
		ID: id,
	}

	err := s.repo.DeleteAutoStore(ctx, deleteID)
	if err != nil {
		s.logger.Errorf("delete autoStore err: %v", err)
		return err
	}
	return nil
}
