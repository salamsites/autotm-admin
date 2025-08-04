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

type BrandHandler struct {
	logger     *slog.Logger
	middleware *shttp.Middleware
	service    repository.BrandService
}

func NewBrandHandler(logger *slog.Logger, middleware *shttp.Middleware, service repository.BrandService) *BrandHandler {
	return &BrandHandler{
		logger:     logger,
		middleware: middleware,
		service:    service,
	}
}

func (h *BrandHandler) BrandRegisterRoutes(r chi.Router) {
	// Body Type
	r.Method("POST", "/create-body-type", h.middleware.Base(h.v1CreateBodyType))
	r.Method("GET", "/get-body-types", h.middleware.Base(h.v1GetBodyTypes))
	r.Method("PUT", "/update-body-type", h.middleware.Base(h.v1UpdateBodyType))
	r.Method("DELETE", "/delete-body-type", h.middleware.Base(h.v1DeleteBodyType))

	// Brand
	r.Method("POST", "/create-brand", h.middleware.Base(h.v1CreateBrand))
	r.Method("GET", "/get-brands", h.middleware.Base(h.v1GetBrands))
	r.Method("PUT", "/update-brand", h.middleware.Base(h.v1UpdateBrand))
	r.Method("DELETE", "/delete-brand", h.middleware.Base(h.v1DeleteBrandCategory))

	// Model
	r.Method("POST", "/create-model", h.middleware.Base(h.v1CreateModel))
	r.Method("GET", "/get-models", h.middleware.Base(h.v1GetModels))
	r.Method("PUT", "/update-model", h.middleware.Base(h.v1UpdateModel))
	r.Method("DELETE", "/delete-model", h.middleware.Base(h.v1DeleteModel))
}

// v1CreateBodyType
// @Summary Create a new body type
// @Description Creates a new body type with the given name, category and image path
// @Tags Body Type
// @Accept json
// @Produce json
// @Param brand body dtos.CreateBodyTypeReq true "Body Type data"
// @Success 200 {object} dtos.ID "Returns created bodyType ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/create-body-type [post]
func (h *BrandHandler) v1CreateBodyType(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var bodyTypeDTO dtos.CreateBodyTypeReq
	errData := json.Unmarshal(body, &bodyTypeDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.CreateBodyType(r.Context(), bodyTypeDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to create body type", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Body Type Create Successfully"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1GetBodyTypes
// @Summary Get body types
// @Description Get paginated list of body types filtered optional search string
// @Tags Body Type
// @Accept json
// @Produce json
// @Param category query string true "Category filter (auto, moto, truck)"
// @Param limit query int false "Limit number of body types to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter body types by name"
// @Success 200 {object} dtos.BodyTypeResult "List of body types with pagination info"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/get-body-types [get]
func (h *BrandHandler) v1GetBodyTypes(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	category := r.URL.Query().Get("category")
	if category == "" {
		result.Message = "category is required"
		return shttp.BadRequest.SetData(result)
	}
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

	brands, err := h.service.GetBodyType(r.Context(), limit, page, category, search)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to get body types", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Body Types Get Successfully"
	result.Data = brands
	return shttp.Success.SetData(result)
}

// v1UpdateBodyType handler
// @Summary Update an existing body type
// @Description Updates body type details by ID
// @Tags Body Type
// @Accept json
// @Produce json
// @Param brand body dtos.UpdateBodyTypeReq true "Body Type data with ID"
// @Success 200 {object} dtos.ID "Returns updated body Type ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/update-body-type [put]
func (h *BrandHandler) v1UpdateBodyType(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var bodyTypeDTO dtos.UpdateBodyTypeReq
	errData := json.Unmarshal(body, &bodyTypeDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body ", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.UpdateBodyType(r.Context(), bodyTypeDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to update body type ", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Body Type Update Successfully"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1DeleteBodyType
// @Summary Delete a body type
// @Description Deletes a body type by ID
// @Tags Body Type
// @Accept json
// @Produce json
// @Param id query int true "Body Type ID to delete"
// @Success 200 {object} string "Body Type deleted successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Body Type not found"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/delete-body-type [delete]
func (h *BrandHandler) v1DeleteBodyType(w http.ResponseWriter, r *http.Request) shttp.Response {
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
		h.logger.Error("invalid body type ID", err)
		return shttp.BadRequest.SetData(result)
	}

	err = h.service.DeleteBodyType(r.Context(), id)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to delete body type", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Body Type deleted successfully"
	return shttp.Success.SetData(result)
}

// v1CreateBrand
// @Summary Create a new brand
// @Description Creates a new brand with the given name and logo path
// @Tags Brand
// @Accept json
// @Produce json
// @Param brand body dtos.CreateBrandReq true "Brand data"
// @Success 200 {object} dtos.ID "Returns created brand ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/create-brand [post]
func (h *BrandHandler) v1CreateBrand(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var brandDTO dtos.CreateBrandReq
	errData := json.Unmarshal(body, &brandDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.CreateBrand(r.Context(), brandDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to create brand", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Brand Create Successfully"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1GetBrands
// @Summary Get brands
// @Description Get paginated list of brands filtered by category and optional search string
// @Tags Brand
// @Accept json
// @Produce json
// @Param category query string true "Category filter (auto, moto, truck)"
// @Param limit query int false "Limit number of brands to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter brands by name"
// @Success 200 {object} dtos.BrandResult "List of brands with pagination info"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/get-brands [get]
func (h *BrandHandler) v1GetBrands(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	category := r.URL.Query().Get("category")
	if category == "" {
		result.Message = "category is required"
		return shttp.BadRequest.SetData(result)
	}
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

	brands, err := h.service.GetBrands(r.Context(), limit, page, category, search)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to get brands", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Get Brands Successfully"
	result.Data = brands
	return shttp.Success.SetData(result)
}

// v1UpdateBrand handler
// @Summary Update an existing brand
// @Description Updates brand details by ID
// @Tags Brand
// @Accept json
// @Produce json
// @Param brand body dtos.UpdateBrandReq true "Brand data with ID"
// @Success 200 {object} dtos.ID "Returns updated brand ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/update-brand [put]
func (h *BrandHandler) v1UpdateBrand(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var brandDTO dtos.UpdateBrandReq
	errData := json.Unmarshal(body, &brandDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.UpdateBrand(r.Context(), brandDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to update brand", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Brand Update Successfully"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1DeleteBrand
// @Summary Delete a brand
// @Description Deletes a brand by ID
// @Tags Brand
// @Accept json
// @Produce json
// @Param id query int true "Brand ID to delete"
// @Param category query string true "Brand Category to delete (auto, moto, truck)"
// @Success 200 {object} string "Brand deleted successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Brand not found"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/delete-brand [delete]
func (h *BrandHandler) v1DeleteBrandCategory(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	category := r.URL.Query().Get("category")
	if category == "" {
		result.Message = "category is required"
		return shttp.BadRequest.SetData(result)
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		result.Message = "id is required"
		return shttp.BadRequest.SetData(result)
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("invalid brand ID", err)
		return shttp.BadRequest.SetData(result)
	}

	err = h.service.DeleteBrandCategory(r.Context(), id, category)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to delete brand", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Brand deleted successfully"
	return shttp.Success.SetData(result)
}

// v1CreateModel
// @Summary Create a new brand model
// @Description Creates a new brand model with the given name and logo path
// @Tags Model
// @Accept json
// @Produce json
// @Param brand body dtos.CreateModelReq true "Model data"
// @Success 200 {object} dtos.ID "Returns created model ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/create-model [post]
func (h *BrandHandler) v1CreateModel(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var modelDTO dtos.CreateModelReq
	errData := json.Unmarshal(body, &modelDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.CreateModel(r.Context(), modelDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to create model", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Model Create Successfully"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1GetModels
// @Summary Get all models
// @Description Get a paginated list of models with optional search
// @Tags Model
// @Accept json
// @Produce json
// @Param category query string true "Category filter (auto, moto, truck)"
// @Param limit query int false "Limit number of models to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter models or brands by name or body types by name"
// @Success 200 {object} dtos.ModelResult "List of models with pagination info"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/get-models [get]
func (h *BrandHandler) v1GetModels(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	category := r.URL.Query().Get("category")
	if category == "" {
		result.Message = "category is required"
		return shttp.BadRequest.SetData(result)
	}
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

	brandModels, err := h.service.GetModels(r.Context(), limit, page, category, search)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to get models", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Models Get Successfully"
	result.Data = brandModels
	return shttp.Success.SetData(result)
}

// v1UpdateModel
// @Summary Update an existing model
// @Description Updates model details by ID
// @Tags Model
// @Accept json
// @Produce json
// @Param brand body dtos.UpdateModelReq true "Model data with ID"
// @Success 200 {object} dtos.ID "Returns updated model ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/update-model [put]
func (h *BrandHandler) v1UpdateModel(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var modelDTO dtos.UpdateModelReq
	errData := json.Unmarshal(body, &modelDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.UpdateModel(r.Context(), modelDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to update model", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Model Update Successfully"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1DeleteModel
// @Summary Delete model
// @Description Deletes model by ID
// @Tags Model
// @Accept json
// @Produce json
// @Param id query int true "Model ID to delete"
// @Success 200 {object} string "Model deleted successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Model not found"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/delete-model [delete]
func (h *BrandHandler) v1DeleteModel(w http.ResponseWriter, r *http.Request) shttp.Response {
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
		h.logger.Error("invalid brand model ID", err)
		return shttp.BadRequest.SetData(result)
	}

	err = h.service.DeleteModel(r.Context(), id)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to delete model", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Model Deleted Successfully"
	return shttp.Success.SetData(result)
}
