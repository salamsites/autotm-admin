package http

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/services/repository"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	shttp "github.com/salamsites/package-http"
	slog "github.com/salamsites/package-log"
	"io"
	"net/http"
	"strconv"
)

type StockHandler struct {
	logger     *slog.Logger
	middleware *shttp.Middleware
	service    repository.StockService
}

func NewStockHandler(logger *slog.Logger, middleware *shttp.Middleware, service repository.StockService) *StockHandler {
	return &StockHandler{
		logger:     logger,
		middleware: middleware,
		service:    service,
	}
}

func (h *StockHandler) StockRegisterRoutes(r chi.Router) {
	r.Method("POST", "/create-stock", h.middleware.Base(h.v1CreateStock))
	//r.Method("GET", "/get-users", h.middleware.Base(h.v1GetUsers))
	r.Method("GET", "/get-stocks", h.middleware.Base(h.v1GetStocks))
	r.Method("PUT", "/update-stock", h.middleware.Base(h.v1UpdateStock))
	r.Method("DELETE", "/delete-stock", h.middleware.Base(h.v1DeleteStock))
}

// v1CreateStock
// @Summary Create a new stock
// @Description Creates a new stock
// @Tags Stock
// @Accept json
// @Produce json
// @Param stock body dtos.CreateStockReq true "Stock data"
// @Success 200 {object} map[string]int64 "Returns created stock ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /stocks/create-stock [post]
func (h *StockHandler) v1CreateStock(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var stock dtos.CreateStockReq
	errData := json.Unmarshal(body, &stock)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.CreateStock(r.Context(), stock)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to create stock", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Created stock Successfully"
	result.Data = map[string]interface{}{
		"id": id,
	}
	return shttp.Success.SetData(result)
}

//// v1GetUsers
//// @Summary Get users
//// @Description Get paginated list of users filtered optional search string
//// @Tags Stock
//// @Accept json
//// @Produce json
//// @Param limit query int false "Limit number of users to return"
//// @Param page query int false "Page number"
//// @Param search query string false "Search string to filter users by name"
//// @Success 200 {object} dtos.GetUserResult "List of users with pagination info successfully"
//// @Failure 400 {object} string "Bad request"
//// @Failure 500 {object} string "Internal server error"
//// @Router /stocks/get-users [get]
//func (h *StockHandler) v1GetUsers(w http.ResponseWriter, r *http.Request) shttp.Response {
//	var result shttp.Result
//	result.Status = false
//
//	limitStr := r.URL.Query().Get("limit")
//	pageStr := r.URL.Query().Get("page")
//	search := r.URL.Query().Get("search")
//
//	limit, err := strconv.ParseInt(limitStr, 10, 64)
//	if err != nil || limit <= 0 {
//		limit = 10
//	}
//	page, err := strconv.ParseInt(pageStr, 10, 64)
//	if err != nil || page <= 0 {
//		page = 1
//	}
//
//	users, err := h.service.GetUsersFromUserService(r.Context(), limit, page, search)
//	if err != nil {
//		result.Message = err.Error()
//		h.logger.Error("unable to get users", err)
//		return shttp.InternalServerError.SetData(result)
//	}
//
//	result.Status = true
//	result.Message = "List of users with pagination info successfully"
//	result.Data = users
//	return shttp.Success.SetData(result)
//}

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
// @Summary Update an existing stock
// @Description Updates stock details by ID
// @Tags Stock
// @Accept json
// @Produce json
// @Param autoStore body dtos.UpdateStockReq true "Stock data with ID"
// @Success 200 {object} dtos.ID "Returns updated stock ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /stocks/update-stock [put]
func (h *StockHandler) v1UpdateStock(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var stockDTO dtos.UpdateStockReq
	errData := json.Unmarshal(body, &stockDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.UpdateStock(r.Context(), stockDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to update stock", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Stock Update Successfully"
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
