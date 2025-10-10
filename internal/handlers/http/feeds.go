package http

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"autotm-admin/internal/services/repository"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Hajymuhammet/elasticsearch-package/filter"
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
	//Cars
	r.Method("GET", "/get-cars", h.middleware.Base(h.v1GetCars))
	r.Method("GET", "/get-car-by-id", h.middleware.Base(h.v1GetCarById))
	r.Method("PUT", "/update-car-status", h.middleware.Base(h.v1UpdateCarStatus))
	r.Method("GET", "/search-cars", h.middleware.Base(h.v1SearchCars))

	//trucks
	r.Method("GET", "/get-trucks", h.middleware.Base(h.v1GetTrucks))
	r.Method("GET", "/get-truck-by-id", h.middleware.Base(h.v1GetTruckById))
	r.Method("PUT", "/update-truck-status", h.middleware.Base(h.v1UpdateTruckStatus))

	//Motors
	r.Method("GET", "/get-motors", h.middleware.Base(h.v1GetMotors))
	r.Method("GET", "/get-moto-by-id", h.middleware.Base(h.v1GetMotoById))
	r.Method("PUT", "/update-moto-status", h.middleware.Base(h.v1UpdateMotoStatus))
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

// v1GetCarById
// @Summary Get car by id
// @Description Get car by ID
// @Tags Cars
// @Accept json
// @Produce json
// @Param id query int true "Car ID to get"
// @Success 200 {object} dtos.Car "Successfully get car by id"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Car not found"
// @Failure 500 {object} string "Internal server error"
// @Router /cars/get-car-by-id [get]
func (h *CarsHandler) v1GetCarById(w http.ResponseWriter, r *http.Request) shttp.Response {
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

	car, err := h.service.GetCarByID(r.Context(), id)
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

	// Validate
	validate := helpers.GetValidator()
	err := validate.Struct(carDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Errorln(err)
		return shttp.BadRequest.SetData(result)
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

// v1GetTrucks
// @Summary Get Trucks
// @Description Get paginated list of trucks filtered optional search string
// @Tags Trucks
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of trucks to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter trucks by name and users by name"
// @Param status query string false "Status string to filter trucks by status (pending, accepted, blocked)"
// @Success 200 {object} dtos.TrucksResp "List of trucks with pagination info successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /cars/get-trucks [get]
func (h *CarsHandler) v1GetTrucks(w http.ResponseWriter, r *http.Request) shttp.Response {
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

	trucks, err := h.service.GetTrucks(r.Context(), limit, page, search, status)
	if err != nil {
		result.Status = false
		result.Message = err.Error()
		h.logger.Error("unable to get trucks", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "List of trucks with pagination info successfully"
	result.Data = trucks
	return shttp.Success.SetData(result)
}

// v1GetTruckById
// @Summary Get truck by id
// @Description Get truck by ID
// @Tags Trucks
// @Accept json
// @Produce json
// @Param id query int true "Truck ID to get"
// @Success 200 {object} dtos.Truck "Successfully get truck by id"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Truck not found"
// @Failure 500 {object} string "Internal server error"
// @Router /cars/get-truck-by-id [get]
func (h *CarsHandler) v1GetTruckById(w http.ResponseWriter, r *http.Request) shttp.Response {
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
		h.logger.Error("invalid truck ID", err)
		return shttp.BadRequest.SetData(result)
	}

	truck, err := h.service.GetTruckByID(r.Context(), id)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to get truck", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Successfully retrieved truck"
	result.Data = truck
	return shttp.Success.SetData(result)
}

// v1UpdateTruckStatus
// @Summary Update Truck Status
// @Description Updates the status of a truck
// @Tags Trucks
// @Accept json
// @Produce json
// @Param Truck body dtos.UpdateTruckStatus true "Truck ID and new Status"
// @Success 200 {object} dtos.ID "Returns updated truck ID"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Truck not found"
// @Failure 500 {object} string "Internal server error"
// @Router /cars/update-truck-status [put]
func (h *CarsHandler) v1UpdateTruckStatus(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var truckDTO dtos.UpdateTruckStatus
	errData := json.Unmarshal(body, &truckDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	// Validate
	validate := helpers.GetValidator()
	err := validate.Struct(truckDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Errorln(err)
		return shttp.BadRequest.SetData(result)
	}

	id, err := h.service.UpdateTruckStatus(r.Context(), truckDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to update truck", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Successfully updated truck"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1GetMotors
// @Summary Get Motors
// @Description Get paginated list of motors filtered optional search string
// @Tags Motors
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of motors to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter motors by name and users by name"
// @Param status query string false "Status string to filter motors by status (pending, accepted, blocked)"
// @Success 200 {object} dtos.MotoResp "List of motors with pagination info successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /cars/get-motors [get]
func (h *CarsHandler) v1GetMotors(w http.ResponseWriter, r *http.Request) shttp.Response {
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

	motors, err := h.service.GetMotors(r.Context(), limit, page, search, status)
	if err != nil {
		result.Status = false
		result.Message = err.Error()
		h.logger.Error("unable to get motors", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "List of motors with pagination info successfully"
	result.Data = motors
	return shttp.Success.SetData(result)
}

// v1GetMotoById
// @Summary Get moto by id
// @Description Get moto by ID
// @Tags Motors
// @Accept json
// @Produce json
// @Param id query int true "Moto ID to get"
// @Success 200 {object} dtos.Moto "Successfully get moto by id"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Moto not found"
// @Failure 500 {object} string "Internal server error"
// @Router /cars/get-moto-by-id [get]
func (h *CarsHandler) v1GetMotoById(w http.ResponseWriter, r *http.Request) shttp.Response {
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
		h.logger.Error("invalid moto ID", err)
		return shttp.BadRequest.SetData(result)
	}

	moto, err := h.service.GetMotoByID(r.Context(), id)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to get moto", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Successfully retrieved moto"
	result.Data = moto
	return shttp.Success.SetData(result)
}

// v1UpdateMotoStatus
// @Summary Update Moto Status
// @Description Updates the status of a moto
// @Tags Motors
// @Accept json
// @Produce json
// @Param Moto body dtos.UpdateMotoStatus true "Moto ID and new Status"
// @Success 200 {object} dtos.ID "Returns updated moto ID"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Car not found"
// @Failure 500 {object} string "Internal server error"
// @Router /cars/update-moto-status [put]
func (h *CarsHandler) v1UpdateMotoStatus(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var motoDTO dtos.UpdateMotoStatus
	errData := json.Unmarshal(body, &motoDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	// Validate
	validate := helpers.GetValidator()
	err := validate.Struct(motoDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Errorln(err)
		return shttp.BadRequest.SetData(result)
	}

	id, err := h.service.UpdateMotoStatus(r.Context(), motoDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to update moto", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Successfully updated moto"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1SearchCars
// @Summary Search Cars
// @Description Search cars with filter params
// @Tags Cars
// @Accept json
// @Produce json
// @Param brand_id query []int64 false "Brand IDs"
// @Param model_id query []int64 false "Model IDs"
// @Param year_min query int false "Minimum year"
// @Param year_max query int false "Maximum year"
// @Param price_min query int false "Minimum price"
// @Param price_max query int false "Maximum price"
// @Param city_id query []int64 false "City IDs"
// @Param engine_type query []string false "Engine types"
// @Param transmission query []string false "Transmission types"
// @Param drive_type query []string false "Drive types"
// @Param body_id query []int64 false "Body IDs"
// @Param mileage_min query int false "Minimum mileage"
// @Param mileage_max query int false "Maximum mileage"
// @Param engine_capacity_min query float64 false "Minimum engine capacity"
// @Param engine_capacity_max query float64 false "Maximum engine capacity"
// @Param color query []string false "Colors"
// @Param is_exchange query bool false "Is exchange"
// @Param is_credit query bool false "Is credit"
// @Param status query []string false "Status"
// @Param created_at_min query string false "Created at min (RFC3339)"
// @Param created_at_max query string false "Created at max (RFC3339)"
// @Success 200 {array} dtos.Car "List of cars filtered"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /cars/search-cars [get]
func (h *CarsHandler) v1SearchCars(w http.ResponseWriter, r *http.Request) shttp.Response {
	q := r.URL.Query()
	var ft filter.CarFilter

	// Helper parse funcs
	parseIntSlice := func(values []string) []int64 {
		var res []int64
		for _, v := range values {
			if id, err := strconv.ParseInt(v, 10, 64); err == nil {
				res = append(res, id)
			}
		}
		return res
	}

	parseFloatPtr := func(value string) *float64 {
		if value == "" {
			return nil
		}
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return &f
		}
		return nil
	}

	parseIntPtr := func(value string) *int64 {
		if value == "" {
			return nil
		}
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return &i
		}
		return nil
	}

	// Fill filter
	ft.BrandID = parseIntSlice(q["brand_id"])
	ft.ModelID = parseIntSlice(q["model_id"])
	ft.YearMin = parseIntPtr(q.Get("year_min"))
	ft.YearMax = parseIntPtr(q.Get("year_max"))
	ft.PriceMin = parseIntPtr(q.Get("price_min"))
	ft.PriceMax = parseIntPtr(q.Get("price_max"))
	ft.CityID = parseIntSlice(q["city_id"])
	ft.EngineType = q["engine_type"]
	ft.Transmission = q["transmission"]
	ft.DriveType = q["drive_type"]
	ft.BodyID = parseIntSlice(q["body_id"])
	ft.MileageMin = parseIntPtr(q.Get("mileage_min"))
	ft.MileageMax = parseIntPtr(q.Get("mileage_max"))
	ft.EngineCapacityMin = parseFloatPtr(q.Get("engine_capacity_min"))
	ft.EngineCapacityMax = parseFloatPtr(q.Get("engine_capacity_max"))
	ft.Color = q["color"]

	if v := q.Get("is_exchange"); v != "" {
		val := v == "true"
		ft.IsExchange = &val
	}
	if v := q.Get("is_credit"); v != "" {
		val := v == "true"
		ft.IsCredit = &val
	}

	ft.Status = q["status"]

	if v := q.Get("created_at_min"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			ft.CreatedAtMin = t
		}
	}
	if v := q.Get("created_at_max"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			ft.CreatedAtMax = t
		}
	}

	// Call service (SearchCars from your ES package)
	cars, err := h.service.SearchCars(r.Context(), &ft)
	if err != nil {
		h.logger.Error("unable to search cars", err)
		return shttp.InternalServerError.SetData(shttp.Result{
			Status:  false,
			Message: err.Error(),
		})
	}

	return shttp.Success.SetData(shttp.Result{
		Status:  true,
		Message: "Successfully searched cars",
		Data:    cars,
	})
}
