package services

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"autotm-admin/internal/models"
	"autotm-admin/internal/repository/storage"
	"context"
	slog "github.com/salamsites/package-log"
)

type StockService struct {
	logger *slog.Logger
	repo   storage.StockRepository
}

func NewStockService(logger *slog.Logger, repo storage.StockRepository) *StockService {
	return &StockService{
		logger: logger,
		repo:   repo,
	}
}

func (s *StockService) CreateStock(ctx context.Context, stock dtos.CreateStockReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(stock); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	newStock := models.Stock{
		UserID:      stock.UserID,
		PhoneNumber: stock.PhoneNumber,
		Email:       stock.Email,
		StoreName:   stock.StoreName,
		Images:      stock.Images,
		Logo:        stock.Logo,
		Address:     stock.Address,
		RegionID:    stock.RegionID,
		CityID:      stock.CityID,
	}

	stockID, err := s.repo.CreateStock(ctx, newStock)
	if err != nil {
		s.logger.Errorf("create err: %v", err)
		return stockID, err
	}
	return stockID, nil
}

//
//func (s *StockService) GetUsersFromUserService(ctx context.Context, limit, page int64, search string) (dtos.GetUserResult, error) {
//	offset := (page - 1) * limit
//	if page <= 0 {
//		page = 1
//		offset = 0
//	}
//
//	users, count, err := s.userService.GetUsers(ctx, limit, offset, search)
//	if err != nil {
//		s.logger.Errorf("get users err: %v", err)
//		return dtos.GetUserResult{}, err
//	}
//
//	result := dtos.GetUserResult{
//		Users: users,
//		Count: count,
//	}
//
//	return result, nil
//}

func (s *StockService) GetStocks(ctx context.Context, limit, page int64, search string) (dtos.StocksResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	stocks, count, err := s.repo.GetStocks(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get autoStores err: %v", err)
		return dtos.StocksResult{}, err
	}
	//
	//userIDMap := make(map[int64]struct{})
	//var id dtos.GetUserByIDsReq
	//for _, a := range autoStores {
	//	if _, exists := userIDMap[a.UserID]; !exists {
	//		userIDMap[a.UserID] = struct{}{}
	//		id.Ids = append(id.Ids, a.UserID)
	//	}
	//}
	//
	//users, err := s.userService.GetUserByIds(ctx, id)
	//if err != nil {
	//	s.logger.Errorf("get user by ids err: %v", err)
	//	return dtos.AutoStoresResult{}, err
	//}
	//
	//userMap := make(map[int64]dtos.GetUsers)
	//for _, user := range users {
	//	userMap[user.Id] = user
	//}

	var dtoStocks []dtos.Stock
	for _, stock := range stocks {
		//user := userMap[autoStore.UserID]
		dtoStocks = append(dtoStocks, dtos.Stock{
			ID:           stock.ID,
			PhoneNumber:  stock.PhoneNumber,
			Email:        stock.Email,
			StoreName:    stock.StoreName,
			Images:       stock.Images,
			Logo:         stock.Logo,
			Address:      stock.Address,
			CityID:       stock.CityID,
			CityNameTM:   stock.CityNameTM,
			CityNameEN:   stock.CityNameEN,
			CityNameRU:   stock.CityNameRU,
			RegionID:     stock.RegionID,
			RegionNameTM: stock.RegionNameTM,
			RegionNameEN: stock.RegionNameEN,
			RegionNameRU: stock.RegionNameRU,
			UserID:       stock.UserID,
			UserName:     stock.UserName,
		})
	}

	result := dtos.StocksResult{
		Stocks: dtoStocks,
		Count:  count,
	}

	return result, nil
}

func (s *StockService) UpdateStock(ctx context.Context, stock dtos.UpdateStockReq) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(stock); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	newStock := models.Stock{
		ID:          stock.ID,
		UserID:      stock.UserID,
		PhoneNumber: stock.PhoneNumber,
		Email:       stock.Email,
		StoreName:   stock.StoreName,
		Images:      stock.Images,
		Logo:        stock.Logo,
		RegionID:    stock.RegionID,
		CityID:      stock.CityID,
		Address:     stock.Address,
	}

	stockID, err := s.repo.UpdateStock(ctx, newStock)
	if err != nil {
		s.logger.Errorf("update stock err: %v", err)
		return id, err
	}
	id.ID = stockID
	return id, nil
}

func (s *StockService) DeleteStock(ctx context.Context, id int64) error {
	deleteID := models.ID{
		ID: id,
	}

	err := s.repo.DeleteStock(ctx, deleteID)
	if err != nil {
		s.logger.Errorf("delete stock err: %v", err)
		return err
	}
	return nil
}
