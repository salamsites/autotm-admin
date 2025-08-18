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

type RegionsHandler struct {
	logger     *slog.Logger
	middleware *shttp.Middleware
	service    repository.RegionsService
}

func NewRegionsHandler(logger *slog.Logger, middleware *shttp.Middleware, service repository.RegionsService) *RegionsHandler {
	return &RegionsHandler{
		logger:     logger,
		middleware: middleware,
		service:    service,
	}
}

func (h *RegionsHandler) RegionsRegisterRoutes(r chi.Router) {
	r.Method("POST", "/create-region", h.middleware.Base(h.v1CreateRegion))
	r.Method("GET", "/get-regions", h.middleware.Base(h.v1GetAllRegions))
	r.Method("PUT", "/update-region", h.middleware.Base(h.v1UpdateRegion))
	r.Method("DELETE", "/delete-region", h.middleware.Base(h.v1DeleteRegion))

	//Cities
	r.Method("POST", "/create-city", h.middleware.Base(h.v1CreateCity))
	r.Method("GET", "/get-cities", h.middleware.Base(h.v1GetAllCities))
	r.Method("PUT", "/update-city", h.middleware.Base(h.v1UpdateCity))
	r.Method("DELETE", "/delete-city", h.middleware.Base(h.v1DeleteCity))
}

// v1CreateRegion
// @Summary Create a new region
// @Description Creates a new region with the given name and region
// @Tags Region
// @Accept json
// @Produce json
// @Param Region body dtos.CreateRegionReq true "Region data"
// @Success 200 {object} dtos.ID "Returns created region ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /regions/create-region [post]
func (h *RegionsHandler) v1CreateRegion(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var regionDTO dtos.CreateRegionReq
	errData := json.Unmarshal(body, &regionDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.CreateRegion(r.Context(), regionDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to create region", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Returns created region ID"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1GetAllRegions
// @Summary Get all regions
// @Description Get a paginated list of regions with optional search
// @Tags Region
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of regions to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter regions by name"
// @Success 200 {object} dtos.RegionResult "List of regions with pagination info Successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /regions/get-regions [get]
func (h *RegionsHandler) v1GetAllRegions(w http.ResponseWriter, r *http.Request) shttp.Response {
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

	regions, err := h.service.GetAllRegions(r.Context(), limit, page, search)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to get regions", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "List of regions with pagination info Successfully"
	result.Data = regions
	return shttp.Success.SetData(result)
}

// v1UpdateRegion handler
// @Summary Update an existing region
// @Description Updates region details by ID
// @Tags Region
// @Accept json
// @Produce json
// @Param Region body dtos.UpdateRegionReq true "Region data with ID"
// @Success 200 {object} dtos.ID "Returns updated region ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /regions/update-region [put]
func (h *RegionsHandler) v1UpdateRegion(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var regionDTO dtos.UpdateRegionReq
	errData := json.Unmarshal(body, &regionDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.UpdateRegion(r.Context(), regionDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to update region", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Returns updated region ID"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1DeleteRegion
// @Summary Delete a region
// @Description Delete a region by ID
// @Tags Region
// @Accept json
// @Produce json
// @Param id query int true "Region ID to delete"
// @Success 200 {object} string "Region deleted successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Brand not found"
// @Failure 500 {object} string "Internal server error"
// @Router /regions/delete-region [delete]
func (h *RegionsHandler) v1DeleteRegion(w http.ResponseWriter, r *http.Request) shttp.Response {
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
		h.logger.Error("invalid region ID", err)
		return shttp.BadRequest.SetData(result)
	}

	err = h.service.DeleteRegion(r.Context(), id)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to delete region", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "region deleted successfully"
	return shttp.Success.SetData(result)
}

// v1CreateCity
// @Summary Create a new city
// @Description Creates a new city with the given name and city
// @Tags City
// @Accept json
// @Produce json
// @Param City body dtos.CreateCityReq true "City data"
// @Success 200 {object} dtos.ID "Returns created city ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /regions/create-city [post]
func (h *RegionsHandler) v1CreateCity(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var cityDTO dtos.CreateCityReq
	errData := json.Unmarshal(body, &cityDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.CreateCity(r.Context(), cityDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to create city", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Returns created city ID"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1GetAllCities
// @Summary Get all cities
// @Description Get a paginated list of cities with optional search
// @Tags City
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of cities to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter cities by name"
// @Param region_id query []int false "Filter cities by multiple region IDs"
// @Success 200 {object} dtos.CityResult "List of cities with pagination info Successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /regions/get-cities [get]
func (h *RegionsHandler) v1GetAllCities(w http.ResponseWriter, r *http.Request) shttp.Response {
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

	var regionIDs []int64
	regionIDStrs := r.URL.Query()["region_id"]
	for _, idStr := range regionIDStrs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil {
			regionIDs = append(regionIDs, id)
		}
	}
	cities, err := h.service.GetAllCities(r.Context(), limit, page, search, regionIDs)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to get cities", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "List of cities with pagination info Successfully"
	result.Data = cities
	return shttp.Success.SetData(result)
}

// v1UpdateCity handler
// @Summary Update an existing city
// @Description Updates city details by ID
// @Tags City
// @Accept json
// @Produce json
// @Param City body dtos.UpdateCityReq true "City data with ID"
// @Success 200 {object} dtos.ID "Returns updated city ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /regions/update-city [put]
func (h *RegionsHandler) v1UpdateCity(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var cityDTO dtos.UpdateCityReq
	errData := json.Unmarshal(body, &cityDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.UpdateCity(r.Context(), cityDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to update city", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Returns updated city ID"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1DeleteCity
// @Summary Delete a city
// @Description Delete a city by ID
// @Tags City
// @Accept json
// @Produce json
// @Param id query int true "City ID to delete"
// @Success 200 {object} string "City deleted successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Brand not found"
// @Failure 500 {object} string "Internal server error"
// @Router /regions/delete-city [delete]
func (h *RegionsHandler) v1DeleteCity(w http.ResponseWriter, r *http.Request) shttp.Response {
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
		h.logger.Error("invalid city ID", err)
		return shttp.BadRequest.SetData(result)
	}

	err = h.service.DeleteCity(r.Context(), id)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to delete city", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "City deleted successfully"
	return shttp.Success.SetData(result)
}
