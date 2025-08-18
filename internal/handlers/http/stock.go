package http

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"autotm-admin/internal/services/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	sminio "github.com/salamsites/minio-pkg"
	"github.com/salamsites/minio-pkg/util"
	shttp "github.com/salamsites/package-http"
	slog "github.com/salamsites/package-log"
)

type StockHandler struct {
	logger           *slog.Logger
	middleware       *shttp.Middleware
	service          repository.StockService
	minioFileClient  sminio.FileClient
	minioImageClient sminio.ImageClient
}

func NewStockHandler(logger *slog.Logger, middleware *shttp.Middleware, service repository.StockService, minioFileClient sminio.FileClient, minioImageClient sminio.ImageClient) *StockHandler {
	return &StockHandler{
		logger:           logger,
		middleware:       middleware,
		service:          service,
		minioFileClient:  minioFileClient,
		minioImageClient: minioImageClient,
	}
}

func (h *StockHandler) StockRegisterRoutes(r chi.Router) {
	r.Method("POST", "/create-stock", h.middleware.Base(h.v1CreateStock))
	r.Method("GET", "/get-stocks", h.middleware.Base(h.v1GetStocks))
	r.Method("PUT", "/update-stock", h.middleware.Base(h.v1UpdateStock))
	r.Method("DELETE", "/delete-stock", h.middleware.Base(h.v1DeleteStock))
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
	phoneNumber := r.FormValue("phone_number")
	email := r.FormValue("email")
	storeName := r.FormValue("store_name")
	regionID := helpers.ParseInt64(r.FormValue("region_id"))
	cityID := helpers.ParseInt64(r.FormValue("city_id"))
	address := r.FormValue("address")

	if userID == 0 || storeName == "" {
		result.Message = "user_id and store_name are required"
		return shttp.BadRequest.SetData(result)
	}

	stock := dtos.CreateStockReq{
		UserID:      userID,
		PhoneNumber: phoneNumber,
		Email:       email,
		StoreName:   storeName,
		RegionID:    regionID,
		CityID:      cityID,
		Address:     address,
	}

	stockID, err := h.service.CreateStock(r.Context(), stock)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to create stock", err)
		return shttp.InternalServerError.SetData(result)
	}

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		h.logger.Error("failed to parse multipart form", err)
		return shttp.BadRequest.SetData("Invalid form data")
	}

	uploadResult, errUpload := h.minioFileClient.UploadFile(r.Context(), r, "image", stockID.ID, helpers.StockImagesSize, util.StockBucket)
	if errUpload.StatusCode != 0 {
		h.logger.Error("failed to upload images", errUpload)
		return shttp.InternalServerError.SetData(errUpload.Message)
	}

	var images []string
	for _, item := range uploadResult.Content {
		if img, ok := item.(util.FeedResultTypeImage); ok {
			images = append(images, img.Path)
		}
	}

	logoFileHeaders := r.MultipartForm.File["logo"]
	var logoPath string
	if len(logoFileHeaders) > 0 {
		logoPath = fmt.Sprintf("%d/logo", stockID.ID)
		errLogo := h.minioImageClient.UploadImage(r.Context(), r, "logo", logoPath, helpers.StockLogoSize, util.FileBucket)
		if errLogo.StatusCode != 0 {
			h.logger.Error("failed to upload logo", errLogo)
			return shttp.InternalServerError.SetData("Failed to upload logo")
		}
	}

	errUpdate := h.service.UpdateStockFiles(r.Context(), stockID, images, logoPath)
	if errUpdate != nil {
		h.logger.Error("unable to update stock files", errUpdate)
		return shttp.InternalServerError.SetData("Failed to update stock files")
	}

	result.Status = true
	result.Message = "Created stock Successfully"
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
// @Success 200 {object} dtos.StocksResult "List of stocks with pagination info successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /stocks/get-stocks [get]
func (h *StockHandler) v1GetStocks(w http.ResponseWriter, r *http.Request) shttp.Response {
	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		limit = 10
	}
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page <= 0 {
		page = 1
	}

	var result shttp.Result

	stocks, err := h.service.GetStocks(r.Context(), limit, page, search)
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

	phoneNumber := r.FormValue("phone_number")
	email := r.FormValue("email")
	regionID := helpers.ParseInt64(r.FormValue("region_id"))
	cityID := helpers.ParseInt64(r.FormValue("city_id"))
	address := r.FormValue("address")

	stockDTO := dtos.UpdateStockReq{
		ID:          stockID,
		UserID:      userID,
		PhoneNumber: phoneNumber,
		Email:       email,
		StoreName:   storeName,
		RegionID:    regionID,
		CityID:      cityID,
		Address:     address,
	}

	id, err := h.service.UpdateStock(r.Context(), stockDTO)
	if err != nil {
		h.logger.Error("unable to update stock", err)
		return shttp.InternalServerError.SetData("Failed to update stock")
	}

	var images []string
	uploadResult, errUpload := h.minioFileClient.UploadFile(r.Context(), r, "image", id.ID, helpers.StockImagesSize, util.StockBucket)
	if errUpload.StatusCode != 0 && errUpload.StatusCode != http.StatusBadRequest {
		h.logger.Error("failed to upload images", errUpload)
		return shttp.InternalServerError.SetData(errUpload.Message)
	}

	if errUpload.StatusCode == 0 {
		for _, item := range uploadResult.Content {
			if img, ok := item.(util.FeedResultTypeImage); ok {
				images = append(images, img.Path)
			}
		}
	}

	logoFileHeaders := r.MultipartForm.File["logo"]
	var logoPath string
	if len(logoFileHeaders) > 0 {
		logoPath = fmt.Sprintf("%d/logo", id.ID)
		errLogo := h.minioImageClient.UploadImage(r.Context(), r, "logo", logoPath, helpers.StockLogoSize, util.FileBucket)
		if errLogo.StatusCode != 0 {
			h.logger.Error("failed to upload logo", errLogo)
			return shttp.InternalServerError.SetData("Failed to upload logo")
		}
	}

	errUpdate := h.service.UpdateStockFiles(r.Context(), id, images, logoPath)
	if errUpdate != nil {
		h.logger.Error("unable to update stock files", errUpdate)
		return shttp.InternalServerError.SetData("Failed to update stock files")
	}

	result.Status = true
	result.Message = "Stock updated successfully"
	result.Data = id
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
