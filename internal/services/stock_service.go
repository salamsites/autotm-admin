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

type StockService struct {
	logger      *slog.Logger
	repo        storage.StockRepository
	userService repository.UserService
	pushService repository.PushService
}

func NewStockService(logger *slog.Logger, repo storage.StockRepository, userService repository.UserService, pushService repository.PushService) *StockService {
	return &StockService{
		logger:      logger,
		repo:        repo,
		userService: userService,
		pushService: pushService,
	}
}

func (s *StockService) CreateStock(ctx context.Context, stock dtos.CreateStockReq) (dtos.ID, error) {
	var id dtos.ID
	validate := helpers.GetValidator()
	if err := validate.Struct(stock); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
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
		Status:      stock.Status,
		Description: stock.Description,
	}

	stockID, err := s.repo.CreateStock(ctx, newStock)
	if err != nil {
		s.logger.Errorf("create err: %v", err)
		return id, err
	}
	id.ID = stockID
	return id, nil
}

func (s *StockService) UpdateStockFiles(ctx context.Context, stockID dtos.ID, images []string, logo string) error {
	if err := s.repo.UpdateStockImages(ctx, stockID.ID, images); err != nil {
		s.logger.Errorf("failed to update stock images: %v", err)
		return err
	}

	if logo != "" {
		if err := s.repo.UpdateStockLogo(ctx, stockID.ID, logo); err != nil {
			s.logger.Errorf("failed to update stock logo: %v", err)
			return err
		}
	}

	return nil
}

func (s *StockService) GetStocks(ctx context.Context, limit, page int64, search, status string) (dtos.StocksResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	stocks, count, err := s.repo.GetStocks(ctx, limit, offset, search, status)
	if err != nil {
		s.logger.Errorf("get autoStores err: %v", err)
		return dtos.StocksResult{}, err
	}

	var dtoStocks []dtos.Stock
	for _, stock := range stocks {
		dtoStocks = append(dtoStocks, dtos.Stock{
			ID:           stock.ID,
			UserID:       stock.UserID,
			UserName:     stock.UserName,
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
			Status:       stock.Status,
			Description:  stock.Description,
		})
	}

	result := dtos.StocksResult{
		Stocks: dtoStocks,
		Count:  count,
	}

	return result, nil
}

func (s *StockService) GetStockByID(ctx context.Context, stockID int64) (dtos.Stock, error) {
	stock, err := s.repo.GetStockByID(ctx, stockID)
	if err != nil {
		s.logger.Errorf("get stock by id err: %v", err)
		return dtos.Stock{}, err
	}

	result := dtos.Stock{
		ID:           stock.ID,
		UserID:       stock.UserID,
		UserName:     stock.UserName,
		PhoneNumber:  stock.PhoneNumber,
		Email:        stock.Email,
		StoreName:    stock.StoreName,
		Images:       stock.Images,
		Logo:         stock.Logo,
		CityID:       stock.CityID,
		CityNameTM:   stock.CityNameTM,
		CityNameEN:   stock.CityNameEN,
		CityNameRU:   stock.CityNameRU,
		RegionID:     stock.RegionID,
		RegionNameTM: stock.RegionNameTM,
		RegionNameEN: stock.RegionNameEN,
		RegionNameRU: stock.RegionNameRU,
		Address:      stock.Address,
		Status:       stock.Status,
		Description:  stock.Description,
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
		Status:      stock.Status,
		Description: stock.Description,
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
	err := s.repo.DeleteStock(ctx, id)
	if err != nil {
		s.logger.Errorf("delete stock err: %v", err)
		return err
	}
	return nil
}

func (s *StockService) UpdateStockStatus(ctx context.Context, stock dtos.UpdateStockStatus) (dtos.ID, error) {
	var id dtos.ID

	validate := helpers.GetValidator()
	if err := validate.Struct(stock); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return id, err
	}

	stockID, err := s.repo.UpdateStockStatus(ctx, stock.ID, stock.Status)
	if err != nil {
		s.logger.Errorf("update stock status err: %v", err)
		return id, err
	}

	userId, err := s.repo.GetUserByStockId(ctx, stockID)
	if err != nil {
		s.logger.Errorf("get GetUserByStockId err: %v", err)
		return id, err
	}

	token, err := s.userService.GetUserFirebaseToken(ctx, userId)
	if err != nil {
		s.logger.Errorf("get GetUserFirebaseToken err: %v", err)
		return id, err
	}

	reqPush := dtos.ReqSendPushDTO{
		Message: stock.Message,
		Token:   token,
	}

	go s.pushService.SendPush(reqPush)

	id.ID = stockID
	return id, nil
}
