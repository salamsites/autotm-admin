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

type AutoStoreHandler struct {
	logger     *slog.Logger
	middleware *shttp.Middleware
	service    repository.AutoStoreService
}

func NewAutoStoreHandler(logger *slog.Logger, middleware *shttp.Middleware, service repository.AutoStoreService) *AutoStoreHandler {
	return &AutoStoreHandler{
		logger:     logger,
		middleware: middleware,
		service:    service,
	}
}

func (h *AutoStoreHandler) AutoStoreRegisterRoutes(r chi.Router) {
	r.Method("POST", "/create-auto-store", h.middleware.Base(h.v1CreateAutoStore))
	r.Method("GET", "/get-users", h.middleware.Base(h.v1GetUsers))
	r.Method("GET", "/get-auto-stores", h.middleware.Base(h.v1GetAutoStores))
	r.Method("PUT", "/update-auto-store", h.middleware.Base(h.v1UpdateAutoStore))
	r.Method("DELETE", "/delete-auto-store", h.middleware.Base(h.v1DeleteAutoStore))
}

// v1CreateAutoStore
// @Summary Create a new auto store
// @Description Creates a new auto store
// @Tags Auto Store
// @Accept json
// @Produce json
// @Param autoStore body dtos.CreateAutoStoreReq true "Auto Store data"
// @Success 200 {object} map[string]int64 "Returns created autoStore ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /auto-store/create-auto-store [post]
func (h *AutoStoreHandler) v1CreateAutoStore(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var autoStore dtos.CreateAutoStoreReq
	errData := json.Unmarshal(body, &autoStore)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.CreateAutoStore(r.Context(), autoStore)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to create auto store", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Created auto Store Successfully"
	result.Data = map[string]interface{}{
		"id": id,
	}
	return shttp.Success.SetData(result)
}

// v1GetUsers
// @Summary Get users
// @Description Get paginated list of users filtered optional search string
// @Tags Auto Store
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of users to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter users by name"
// @Success 200 {object} dtos.GetUserResult "List of users with pagination info successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /auto-store/get-users [get]
func (h *AutoStoreHandler) v1GetUsers(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

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

	users, err := h.service.GetUsersFromUserService(r.Context(), limit, page, search)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to get users", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "List of users with pagination info successfully"
	result.Data = users
	return shttp.Success.SetData(result)
}

// v1GetAutoStores
// @Summary Get AutoStores
// @Description Get paginated list of auto stores filtered optional search string
// @Tags Auto Store
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of users to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter auto stores by name"
// @Success 200 {object} dtos.AutoStoresResult "List of auto stores with pagination info successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /auto-store/get-auto-stores [get]
func (h *AutoStoreHandler) v1GetAutoStores(w http.ResponseWriter, r *http.Request) shttp.Response {
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

	autoStores, err := h.service.GetAutoStores(r.Context(), limit, page, search)
	if err != nil {
		result.Status = false
		result.Message = err.Error()
		h.logger.Error("unable to get autoStores", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "List of auto stores with pagination info successfully"
	result.Data = autoStores
	return shttp.Success.SetData(result)
}

// v1UpdateAutoStore
// @Summary Update an existing auto store
// @Description Updates auto store details by ID
// @Tags Auto Store
// @Accept json
// @Produce json
// @Param autoStore body dtos.UpdateAutoStoreReq true "AutoStore data with ID"
// @Success 200 {object} dtos.ID "Returns updated autoStore ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /auto-store/update-auto-store [put]
func (h *AutoStoreHandler) v1UpdateAutoStore(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var autoStoreDTO dtos.UpdateAutoStoreReq
	errData := json.Unmarshal(body, &autoStoreDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.UpdateAutoStore(r.Context(), autoStoreDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to update auto store", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Auto Store Update Successfully"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1DeleteAutoStore
// @Summary Delete Auto Store
// @Description Deletes auto store by ID
// @Tags Auto Store
// @Accept json
// @Produce json
// @Param id query int true "Model ID to delete"
// @Success 200 {object} string "Model deleted successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Model not found"
// @Failure 500 {object} string "Internal server error"
// @Router /auto-store/delet-auto-store [delete]
func (h *AutoStoreHandler) v1DeleteAutoStore(w http.ResponseWriter, r *http.Request) shttp.Response {
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
		h.logger.Error("invalid auto Store ID", err)
		return shttp.BadRequest.SetData(result)
	}

	err = h.service.DeleteAutoStore(r.Context(), id)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to delete auto store", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Auto Store Deleted Successfully"
	return shttp.Success.SetData(result)
}
