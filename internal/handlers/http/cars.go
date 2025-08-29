package http

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/services/repository"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	shttp "github.com/salamsites/package-http"
	slog "github.com/salamsites/package-log"
)

type CarsHandler struct {
	logger     *slog.Logger
	middleware *shttp.Middleware
	service    repository.CarsService
}

func NewCarsHandler(logger *slog.Logger, middleware *shttp.Middleware, service repository.CarsService) *CarsHandler {
	return &CarsHandler{
		logger:     logger,
		middleware: middleware,
		service:    service,
	}
}

func (h *CarsHandler) CarsRegisterRoutes(r chi.Router) {
	r.Method("GET", "/get-cars", h.middleware.Base(h.v1GetCars))
	r.Method("GET", "/get-car-by-id", h.middleware.Base(h.v1GetCarsById))
	r.Method("PUT", "/update-car-status", h.middleware.Base(h.v1UpdateCarStatus))
}

// v1GetCars
// @Summary Get Cars
// @Description Get paginated list of cars filtered optional search string
// @Tags Cars
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of cars to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter cars by name and users by name"
// @Param status query string false "Status string to filter cars by status (pending, accepted, blocked)"
// @Success 200 {object} dtos.CarsResp "List of cars with pagination info successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /cars/get-cars [get]
func (h *CarsHandler) v1GetCars(w http.ResponseWriter, r *http.Request) shttp.Response {
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

	cars, err := h.service.GetCars(r.Context(), limit, page, search, status)
	if err != nil {
		result.Status = false
		result.Message = err.Error()
		h.logger.Error("unable to get cars", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "List of cars with pagination info successfully"
	result.Data = cars
	return shttp.Success.SetData(result)
}

// v1GetCarsById
// @Summary Get car by id
// @Description Get car by ID
// @Tags Cars
// @Accept json
// @Produce json
// @Param id query int true "Car ID to get"
// @Success 200 {object} dtos.Car "Successfully get car by id"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Brand not found"
// @Failure 500 {object} string "Internal server error"
// @Router /cars/get-car-by-id [get]
func (h *CarsHandler) v1GetCarsById(w http.ResponseWriter, r *http.Request) shttp.Response {
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

	car, err := h.service.GetCarsByID(r.Context(), id)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to get car", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Successfully retrieved car"
	result.Data = car
	return shttp.Success.SetData(result)
}

// v1UpdateCarStatus
// @Summary Update Car Status
// @Description Updates the status of a car
// @Tags Cars
// @Accept json
// @Produce json
// @Param Car body dtos.UpdateCarStatus true "Car ID and new Status"
// @Success 200 {object} dtos.ID "Returns updated car ID"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Car not found"
// @Failure 500 {object} string "Internal server error"
// @Router /cars/update-car-status [put]
func (h *CarsHandler) v1UpdateCarStatus(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var carDTO dtos.UpdateCarStatus
	errData := json.Unmarshal(body, &carDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.UpdateCarStatus(r.Context(), carDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to update car", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Successfully updated car"
	result.Data = id
	return shttp.Success.SetData(result)
}
