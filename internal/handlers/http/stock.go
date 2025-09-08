package http

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"autotm-admin/internal/services/repository"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Hajymuhammet/elasticsearch-package/pkg/querybuilder"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	sminio "github.com/salamsites/minio-pkg"
	"github.com/salamsites/minio-pkg/util"
	shttp "github.com/salamsites/package-http"
	slog "github.com/salamsites/package-log"
	spsql "github.com/salamsites/package-psql"
)

type StockHandler struct {
	logger           *slog.Logger
	middleware       *shttp.Middleware
	clientPsql       spsql.Client
	service          repository.StockService
	minioFileClient  sminio.FileClient
	minioImageClient sminio.ImageClient
}

func NewStockHandler(logger *slog.Logger, middleware *shttp.Middleware, clientPsql spsql.Client, service repository.StockService, minioFileClient sminio.FileClient, minioImageClient sminio.ImageClient) *StockHandler {
	return &StockHandler{
		logger:           logger,
		middleware:       middleware,
		clientPsql:       clientPsql,
		service:          service,
		minioFileClient:  minioFileClient,
		minioImageClient: minioImageClient,
	}
}

func (h *StockHandler) StockRegisterRoutes(r chi.Router) {
	r.Method("POST", "/create-stock", h.middleware.Base(h.v1CreateStock))
	r.Method("GET", "/get-stocks", h.middleware.Base(h.v1GetStocks))
	r.Method("GET", "/get-stock-by-id", h.middleware.Base(h.v1GetStockByID))
	r.Method("PUT", "/update-stock", h.middleware.Base(h.v1UpdateStock))
	r.Method("DELETE", "/delete-stock", h.middleware.Base(h.v1DeleteStock))
	r.Method("PUT", "/update-stock-status", h.middleware.Base(h.v1UpdateStockStatus))
	r.Method("GET", "/search-stocks", h.middleware.Base(h.v1SearchStocks))
	r.Method("POST", "/search-stocks-post", h.middleware.Base(h.v1SearchStocksPost))
}

// v1CreateStock
// @Summary Create a new stock with images and logo upload
// @Description Creates a new stock entry and uploads multiple images and optional logo linked to the stock
// @Tags Stock
// @Accept multipart/form-data
// @Produce json
// @Param user_id formData int true "User ID"
// @Param phone_number formData string false "Phone number"
// @Param email formData string false "Email"
// @Param store_name formData string true "Store name"
// @Param region_id formData int false "Region ID"
// @Param city_id formData int false "City ID"
// @Param address formData string false "Address"
// @Param image formData []file true "Image file(s)"
// @Param logo formData file false "Logo image file"
// @Success 200 {object} dtos.ID "Returns created stock ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /stocks/create-stock [post]
func (h *StockHandler) v1CreateStock(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	userID := helpers.ParseInt64(r.FormValue("user_id"))
	storeName := r.FormValue("store_name")
	if userID == 0 || storeName == "" {
		result.Message = "user_id and store_name are required"
		return shttp.BadRequest.SetData(result)
	}

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		h.logger.Error("failed to parse multipart form", err)
		return shttp.BadRequest.SetData("Invalid form data")
	}

	stock := dtos.CreateStockReq{
		UserID:      userID,
		PhoneNumber: r.FormValue("phone_number"),
		Email:       r.FormValue("email"),
		StoreName:   storeName,
		RegionID:    helpers.ParseInt64(r.FormValue("region_id")),
		CityID:      helpers.ParseInt64(r.FormValue("city_id")),
		Address:     r.FormValue("address"),
		Status:      "accepted", // initial status
	}

	trManager := manager.Must(trmpgx.NewDefaultFactory(h.clientPsql.Pool()))

	var stockID dtos.ID
	err := trManager.Do(r.Context(), func(ctx context.Context) error {
		var err error
		stockID, err = h.service.CreateStock(ctx, stock)
		return err
	})
	if err != nil {
		result.Message = fmt.Sprintf("unable to create stock: %v", err)
		h.logger.Error("create stock transaction failed", err)
		return shttp.InternalServerError.SetData(result)
	}

	go func(stockID dtos.ID) {
		ctx := context.Background() // independent context
		uploadResult, errUpload := h.minioFileClient.UploadFile(ctx, r, "image", stockID.ID, helpers.StockImagesSize, util.StockBucket)
		if errUpload.StatusCode != 0 {
			h.logger.Error("failed to upload images", errUpload)
			return
		}

		logoFileHeaders := r.MultipartForm.File["logo"]
		if len(logoFileHeaders) > 0 {
			logoPath := fmt.Sprintf("%d/logo", stockID.ID)
			errLogo := h.minioImageClient.UploadImage(ctx, r, "logo", logoPath, helpers.StockLogoSize, util.StockBucket)
			if errLogo.StatusCode != 0 {
				h.logger.Error("failed to upload logo", errLogo)
				return
			}
		}

		if errUpdate := h.service.UpdateStockFiles(ctx, stockID, uploadResult, helpers.StockLogoSize); errUpdate != nil {
			h.logger.Error("unable to update stock files", errUpdate)
		}
	}(stockID)

	result.Status = true
	result.Message = "Created stock successfully, files are uploading"
	result.Data = stockID
	return shttp.Success.SetData(result)
}

// v1GetStocks
// @Summary Get Stocks
// @Description Get paginated list of stocks filtered optional search string
// @Tags Stock
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of stocks to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter stocks by name and users by name"
// @Param status query string false "Status string to filter stocks by status (waiting, accepted, blocked)"
// @Success 200 {object} dtos.StocksResult "List of stocks with pagination info successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /stocks/get-stocks [get]
func (h *StockHandler) v1GetStocks(w http.ResponseWriter, r *http.Request) shttp.Response {
	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		limit = 10
	}
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page <= 0 {
		page = 1
	}

	var result shttp.Result

	stocks, err := h.service.GetStocks(r.Context(), limit, page, search, status)
	if err != nil {
		result.Status = false
		result.Message = err.Error()
		h.logger.Error("unable to get stocks", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "List of stocks with pagination info successfully"
	result.Data = stocks
	return shttp.Success.SetData(result)
}

// v1GetStockByID
// @Summary Get stock by id
// @Description Get stock by ID
// @Tags Stock
// @Accept json
// @Produce json
// @Param id query int true "Stock ID to get"
// @Success 200 {object} dtos.Stock "Successfully get stock by id"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Brand not found"
// @Failure 500 {object} string "Internal server error"
// @Router /stocks/get-stock-by-id [get]
func (h *StockHandler) v1GetStockByID(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		result.Message = "id is required"
		return shttp.BadRequest.SetData(result)
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("invalid car ID", err)
		return shttp.BadRequest.SetData(result)
	}

	stock, err := h.service.GetStockByID(r.Context(), id)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to get stock", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Successfully retrieved stock"
	result.Data = stock
	return shttp.Success.SetData(result)
}

// v1UpdateStock
// @Summary Update an existing stock with images and logo upload
// @Description Updates stock details and uploads multiple images and optional logo linked to the stock
// @Tags Stock
// @Accept multipart/form-data
// @Produce json
// @Param id formData int true "Stock ID"
// @Param user_id formData int true "User ID"
// @Param phone_number formData string false "Phone number"
// @Param email formData string false "Email"
// @Param store_name formData string true "Store name"
// @Param region_id formData int false "Region ID"
// @Param city_id formData int false "City ID"
// @Param address formData string false "Address"
// @Param image formData []file false "Image file(s)"
// @Param logo formData file false "Logo image file"
// @Param status formData string false "Status (waiting, accepted, blocked)"
// @Success 200 {object} dtos.ID "Returns updated stock ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /stocks/update-stock [put]
func (h *StockHandler) v1UpdateStock(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		h.logger.Error("failed to parse multipart form", err)
		return shttp.BadRequest.SetData("Invalid form data")
	}

	stockIDStr := r.FormValue("id")
	stockID, err := strconv.ParseInt(stockIDStr, 10, 64)
	if err != nil || stockID == 0 {
		return shttp.BadRequest.SetData("Invalid or missing stock ID")
	}

	userID := helpers.ParseInt64(r.FormValue("user_id"))
	if userID == 0 {
		return shttp.BadRequest.SetData("user_id is required")
	}

	storeName := r.FormValue("store_name")
	if storeName == "" {
		return shttp.BadRequest.SetData("store_name is required")
	}

	stockDTO := dtos.UpdateStockReq{
		ID:          stockID,
		UserID:      userID,
		PhoneNumber: r.FormValue("phone_number"),
		Email:       r.FormValue("email"),
		StoreName:   storeName,
		RegionID:    helpers.ParseInt64(r.FormValue("region_id")),
		CityID:      helpers.ParseInt64(r.FormValue("city_id")),
		Address:     r.FormValue("address"),
		Status:      r.FormValue("status"),
	}

	trManager := manager.Must(trmpgx.NewDefaultFactory(h.clientPsql.Pool()))
	var updatedID dtos.ID
	err = trManager.Do(r.Context(), func(ctx context.Context) error {
		var err error
		updatedID, err = h.service.UpdateStock(ctx, stockDTO)
		return err
	})
	if err != nil {
		result.Message = fmt.Sprintf("unable to update stock: %v", err)
		h.logger.Error("update stock transaction failed", err)
		return shttp.InternalServerError.SetData(result)
	}

	go func(id dtos.ID) {
		ctx := context.Background() // independent context
		uploadResult, errUpload := h.minioFileClient.UploadFile(ctx, r, "image", id.ID, helpers.StockImagesSize, util.StockBucket)
		if errUpload.StatusCode != 0 && errUpload.StatusCode != http.StatusBadRequest {
			h.logger.Error("failed to upload images", errUpload)
			return
		}

		logoFileHeaders := r.MultipartForm.File["logo"]
		if len(logoFileHeaders) > 0 {
			logoPath := fmt.Sprintf("%d/logo", id.ID)
			errLogo := h.minioImageClient.UploadImage(ctx, r, "logo", logoPath, helpers.StockLogoSize, util.FileBucket)
			if errLogo.StatusCode != 0 {
				h.logger.Error("failed to upload logo", errLogo)
				return
			}
		}

		if errUpdate := h.service.UpdateStockFiles(ctx, id, uploadResult, helpers.StockLogoSize); errUpdate != nil {
			h.logger.Error("unable to update stock files", errUpdate)
		}
	}(updatedID)

	result.Status = true
	result.Message = "Stock updated successfully"
	result.Data = updatedID
	return shttp.Success.SetData(result)
}

// v1DeleteStock
// @Summary Delete Stock
// @Description Deletes stock by ID
// @Tags Stock
// @Accept json
// @Produce json
// @Param id query int true "Stock ID to delete"
// @Success 200 {object} string "Stock deleted successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Stock not found"
// @Failure 500 {object} string "Internal server error"
// @Router /stocks/delete-stock [delete]
func (h *StockHandler) v1DeleteStock(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		result.Message = "id is required"
		return shttp.BadRequest.SetData(result)
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("invalid stock ID", err)
		return shttp.BadRequest.SetData(result)
	}

	err = h.service.DeleteStock(r.Context(), id)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to delete auto store", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Stock Deleted Successfully"
	return shttp.Success.SetData(result)
}

// v1UpdateStockStatus
// @Summary Update Stock Status
// @Description Updates the status of a stock
// @Tags Stock
// @Accept json
// @Produce json
// @Param Stock body dtos.UpdateStockStatus true "Stock ID and new Status"
// @Success 200 {object} dtos.ID "Returns updated stock ID"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Stock not found"
// @Failure 500 {object} string "Internal server error"
// @Router /stocks/update-stock-status [put]
func (h *StockHandler) v1UpdateStockStatus(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var stockDTO dtos.UpdateStockStatus
	errData := json.Unmarshal(body, &stockDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.UpdateStockStatus(r.Context(), stockDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to update stock", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Successfully updated stock"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1SearchStocks
// @Summary Search Stocks with Elasticsearch
// @Description Search stocks using Elasticsearch with advanced query capabilities
// @Tags Stock
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param status query string false "Status filter"
// @Param region_id query int false "Region ID filter"
// @Param city_id query int false "City ID filter"
// @Param from query int false "From index for pagination"
// @Param size query int false "Size of results for pagination"
// @Param sort query string false "Sort field"
// @Success 200 {object} dtos.StocksResult "Search results"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /stocks/search-stocks [get]
func (h *StockHandler) v1SearchStocks(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	search := r.URL.Query().Get("q")
	status := r.URL.Query().Get("status")
	regionID := helpers.ParseInt64(r.URL.Query().Get("region_id"))
	cityID := helpers.ParseInt64(r.URL.Query().Get("city_id"))

	// Pagination
	from, _ := strconv.Atoi(r.URL.Query().Get("from"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if size == 0 {
		size = 10
	}
	if from < 0 {
		from = 0
	}

	// Sort
	var sort []map[string]any
	sortField := r.URL.Query().Get("sort")
	if sortField != "" {
		sort = append(sort, map[string]any{sortField: "desc"})
	}

	// Build Elasticsearch query using QueryBuilder
	qb := querybuilder.New()

	// Build filters
	var filters []map[string]any

	if status != "" {
		filters = append(filters, map[string]any{
			"term": map[string]any{"status": status},
		})
	}

	if regionID > 0 {
		filters = append(filters, map[string]any{
			"term": map[string]any{"region_id": regionID},
		})
	}

	if cityID > 0 {
		filters = append(filters, map[string]any{
			"term": map[string]any{"city_id": cityID},
		})
	}

	if search != "" {
		// Use multi_match for search across multiple fields
		searchQuery := map[string]any{
			"multi_match": map[string]any{
				"query": search,
				"fields": []string{
					"store_name", "user_name", "description", "address",
					"city_name_tm", "city_name_en", "city_name_ru",
					"region_name_tm", "region_name_en", "region_name_ru",
				},
			},
		}

		// If we have both search and filters, combine them with bool query
		if len(filters) > 0 {
			qb.Bool(
				[]map[string]any{searchQuery}, // must (search)
				nil,                           // should
				nil,                           // must_not
				filters,                       // filter
			)
		} else {
			// Only search, no filters
			qb.MultiMatch([]string{
				"store_name", "user_name", "description", "address",
				"city_name_tm", "city_name_en", "city_name_ru",
				"region_name_tm", "region_name_en", "region_name_ru",
			}, search)
		}
	} else if len(filters) > 0 {
		// Only filters, no search
		qb.Bool(
			nil,     // must
			nil,     // should
			nil,     // must_not
			filters, // filter
		)
	} else {
		qb.Bool(nil, nil, nil, nil)
	}

	query := qb.Build()

	stocks, total, err := h.service.SearchStocksES(r.Context(), query, from, size, sort)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to search stocks", err)
		return shttp.InternalServerError.SetData(result)
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

	result.Status = true
	result.Message = "Search completed successfully"
	result.Data = dtos.StocksResult{
		Stocks: dtoStocks,
		Count:  total,
	}
	return shttp.Success.SetData(result)
}

// v1SearchStocksPost
// @Summary Search Stocks with Elasticsearch (POST)
// @Description Search stocks using Elasticsearch with advanced query capabilities via POST
// @Tags Stock
// @Accept json
// @Produce json
// @Param search body map[string]any true "Elasticsearch query DSL"
// @Success 200 {object} dtos.StocksResult "Search results"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /stocks/search-stocks-post [post]
func (h *StockHandler) v1SearchStocksPost(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, err := io.ReadAll(r.Body)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to read request body", err)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var query map[string]any
	if err := json.Unmarshal(body, &query); err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to unmarshal query", err)
		return shttp.BadRequest.SetData(result)
	}

	from, _ := strconv.Atoi(r.URL.Query().Get("from"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if size == 0 {
		size = 10
	}
	if from < 0 {
		from = 0
	}

	// Sort
	var sort []map[string]any
	sortField := r.URL.Query().Get("sort")
	if sortField != "" {
		sort = append(sort, map[string]any{sortField: "desc"})
	}

	stocks, total, err := h.service.SearchStocksES(r.Context(), query, from, size, sort)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to search stocks", err)
		return shttp.InternalServerError.SetData(result)
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

	result.Status = true
	result.Message = "Search completed successfully"
	result.Data = dtos.StocksResult{
		Stocks: dtoStocks,
		Count:  total,
	}
	return shttp.Success.SetData(result)
}
