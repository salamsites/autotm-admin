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
// @Param Region body dtos.Region true "Region data"
// @Success 200 {object} map[string]int64 "Returns created region ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /regions/create-region [post]
func (h *RegionsHandler) v1CreateRegion(w http.ResponseWriter, r *http.Request) shttp.Response {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(errBody.Error())
	}
	defer r.Body.Close()

	var regionDTO dtos.Region
	errData := json.Unmarshal(body, &regionDTO)
	if errData != nil {
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(errData.Error())
	}

	id, err := h.service.CreateRegion(r.Context(), regionDTO)
	if err != nil {
		h.logger.Error("unable to create region", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(map[string]interface{}{
		"id": id,
	})
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
// @Success 200 {object} dtos.RegionResult "List of regions with pagination info"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /regions/get-regions [get]
func (h *RegionsHandler) v1GetAllRegions(w http.ResponseWriter, r *http.Request) shttp.Response {
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

	roles, err := h.service.GetAllRegions(r.Context(), limit, page, search)
	if err != nil {
		h.logger.Error("unable to get regions", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(roles)
}

// v1UpdateRegion handler
// @Summary Update an existing region
// @Description Updates region details by ID
// @Tags Region
// @Accept json
// @Produce json
// @Param Region body dtos.Region true "Region data with ID"
// @Success 200 {object} map[string]int64 "Returns updated region ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /regions/update-region [put]
func (h *RegionsHandler) v1UpdateRegion(w http.ResponseWriter, r *http.Request) shttp.Response {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(errBody.Error())
	}
	defer r.Body.Close()

	var regionDTO dtos.Region
	errData := json.Unmarshal(body, &regionDTO)
	if errData != nil {
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(errData.Error())
	}

	id, err := h.service.UpdateRegion(r.Context(), regionDTO)
	if err != nil {
		h.logger.Error("unable to update region", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(map[string]interface{}{
		"id": id,
	})
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
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return shttp.BadRequest.SetData("missing region ID")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("invalid region ID", err)
		return shttp.BadRequest.SetData("invalid region ID")
	}

	err = h.service.DeleteRegion(r.Context(), id)
	if err != nil {
		h.logger.Error("unable to delete region", err)
		return shttp.InternalServerError.SetData(err.Error())
	}

	return shttp.Success.SetData("region deleted successfully")
}

// v1CreateCity
// @Summary Create a new city
// @Description Creates a new city with the given name and city
// @Tags City
// @Accept json
// @Produce json
// @Param City body dtos.City true "City data"
// @Success 200 {object} map[string]int64 "Returns created city ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /regions/create-city [post]
func (h *RegionsHandler) v1CreateCity(w http.ResponseWriter, r *http.Request) shttp.Response {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(errBody.Error())
	}
	defer r.Body.Close()

	var cityDTO dtos.City
	errData := json.Unmarshal(body, &cityDTO)
	if errData != nil {
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(errData.Error())
	}

	id, err := h.service.CreateCity(r.Context(), cityDTO)
	if err != nil {
		h.logger.Error("unable to create city", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(map[string]interface{}{
		"id": id,
	})
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
// @Success 200 {object} dtos.CityResult "List of cities with pagination info"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /regions/get-cities [get]
func (h *RegionsHandler) v1GetAllCities(w http.ResponseWriter, r *http.Request) shttp.Response {
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

	roles, err := h.service.GetAllCities(r.Context(), limit, page, search)
	if err != nil {
		h.logger.Error("unable to get cities", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(roles)
}

// v1UpdateCity handler
// @Summary Update an existing city
// @Description Updates city details by ID
// @Tags City
// @Accept json
// @Produce json
// @Param City body dtos.City true "City data with ID"
// @Success 200 {object} map[string]int64 "Returns updated city ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /regions/update-city [put]
func (h *RegionsHandler) v1UpdateCity(w http.ResponseWriter, r *http.Request) shttp.Response {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(errBody.Error())
	}
	defer r.Body.Close()

	var cityDTO dtos.City
	errData := json.Unmarshal(body, &cityDTO)
	if errData != nil {
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(errData.Error())
	}

	id, err := h.service.UpdateCity(r.Context(), cityDTO)
	if err != nil {
		h.logger.Error("unable to update city", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(map[string]interface{}{
		"id": id,
	})
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
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return shttp.BadRequest.SetData("missing city ID")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("invalid city ID", err)
		return shttp.BadRequest.SetData("invalid city ID")
	}

	err = h.service.DeleteCity(r.Context(), id)
	if err != nil {
		h.logger.Error("unable to delete city", err)
		return shttp.InternalServerError.SetData(err.Error())
	}

	return shttp.Success.SetData("city deleted successfully")
}
