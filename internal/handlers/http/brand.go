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
	r.Method("POST", "/upload-image", h.middleware.Base(h.v1UploadImage))
	r.Method("POST", "/create-brand", h.middleware.Base(h.v1CreateBrand))
	r.Method("GET", "/get-brands", h.middleware.Base(h.v1GetBrands))
	r.Method("PUT", "/update-brand", h.middleware.Base(h.v1UpdateBrand))
	r.Method("DELETE", "/delete-brand", h.middleware.Base(h.v1DeleteBrand))
	r.Method("POST", "/create-model", h.middleware.Base(h.v1CreateModel))
	r.Method("GET", "/get-models", h.middleware.Base(h.v1GetModels))
	r.Method("PUT", "/update-model", h.middleware.Base(h.v1UpdateModel))
	r.Method("DELETE", "/delete-model", h.middleware.Base(h.v1DeleteModel))
}

// v1UploadImage
// @Summary Upload an image
// @Description Uploads an image file and returns the file path
// @Tags Upload Image
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file to upload"
// @Success 200 "Returns the uploaded image path Successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/upload-image [post]
func (h *BrandHandler) v1UploadImage(w http.ResponseWriter, r *http.Request) shttp.Response {
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		h.logger.Error("unable to get uploaded file", err)
		return shttp.BadRequest.SetData(err.Error())
	}
	defer file.Close()

	imagePath, err := h.service.UploadImage(file, fileHeader)
	if err != nil {
		h.logger.Error("unable to upload image error", err)
		return shttp.InternalServerError.SetData(err.Error())
	}

	return shttp.Success.SetData(imagePath)
}

// v1CreateBrand
// @Summary Create a new brand
// @Description Creates a new brand with the given name and logo path
// @Tags Brand
// @Accept json
// @Produce json
// @Param brand body dtos.V1BrandDTO true "Brand data"
// @Success 200 {object} map[string]int64 "Returns created brand ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/create-brand [post]
func (h *BrandHandler) v1CreateBrand(w http.ResponseWriter, r *http.Request) shttp.Response {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(errBody.Error())
	}
	defer r.Body.Close()

	var brandDTO dtos.V1BrandDTO
	errData := json.Unmarshal(body, &brandDTO)
	if errData != nil {
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(errData.Error())
	}

	id, err := h.service.CreateBrand(r.Context(), brandDTO)
	if err != nil {
		h.logger.Error("unable to create brand", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(map[string]interface{}{
		"id": id,
	})
}

// v1GetBrands
// @Summary Get all brands
// @Description Get a paginated list of brands with optional search
// @Tags Brand
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of brands to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter brands by name"
// @Success 200 {object} models.BrandResult "List of brands with pagination info"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/get-brands [get]
func (h *BrandHandler) v1GetBrands(w http.ResponseWriter, r *http.Request) shttp.Response {
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

	brands, err := h.service.GetBrands(r.Context(), limit, page, search)
	if err != nil {
		h.logger.Error("unable to get brands", err, err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(brands)
}

// v1UpdateBrand handler
// @Summary Update an existing brand
// @Description Updates brand details by ID
// @Tags Brand
// @Accept json
// @Produce json
// @Param brand body dtos.V1BrandDTO true "Brand data with ID"
// @Success 200 {object} map[string]int64 "Returns updated brand ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/update-brand [put]
func (h *BrandHandler) v1UpdateBrand(w http.ResponseWriter, r *http.Request) shttp.Response {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(errBody.Error())
	}
	defer r.Body.Close()

	var brandDTO dtos.V1BrandDTO
	errData := json.Unmarshal(body, &brandDTO)
	if errData != nil {
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(errData.Error())
	}

	id, err := h.service.UpdateBrand(r.Context(), brandDTO)
	if err != nil {
		h.logger.Error("unable to update brand", err, err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(map[string]interface{}{
		"id": id,
	})
}

// v1DeleteBrand
// @Summary Delete a brand
// @Description Deletes a brand by ID
// @Tags Brand
// @Accept json
// @Produce json
// @Param id query int true "Brand ID to delete"
// @Success 200 {object} string "Brand deleted successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Brand not found"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/delete-brand [delete]
func (h *BrandHandler) v1DeleteBrand(w http.ResponseWriter, r *http.Request) shttp.Response {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return shttp.BadRequest.SetData("missing brand ID")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("invalid brand ID", err)
		return shttp.BadRequest.SetData("invalid brand ID")
	}

	err = h.service.DeleteBrand(r.Context(), id)
	if err != nil {
		h.logger.Error("unable to delete brand", err)
		return shttp.InternalServerError.SetData(err.Error())
	}

	return shttp.Success.SetData("brand deleted successfully")
}

// v1CreateBrandModel
// @Summary Create a new brand model
// @Description Creates a new brand model with the given name and logo path
// @Tags Brand Model
// @Accept json
// @Produce json
// @Param brand body dtos.V1BrandModelDTO true "Brand Model data"
// @Success 200 {object} map[string]int64 "Returns created model ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/create-model [post]
func (h *BrandHandler) v1CreateModel(w http.ResponseWriter, r *http.Request) shttp.Response {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(errBody.Error())
	}
	defer r.Body.Close()

	var modelDTO dtos.V1BrandModelDTO
	errData := json.Unmarshal(body, &modelDTO)
	if errData != nil {
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(errData.Error())
	}

	id, err := h.service.CreateBrandModel(r.Context(), modelDTO)
	if err != nil {
		h.logger.Error("unable to create brand model", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(map[string]interface{}{
		"id": id,
	})
}

// v1GetBrandModels
// @Summary Get all brand models
// @Description Get a paginated list of brand models with optional search
// @Tags Brand Model
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of brand models to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter brand models or brands by name"
// @Success 200 {object} models.BrandModelResult "List of brand models with pagination info"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/get-models [get]
func (h *BrandHandler) v1GetModels(w http.ResponseWriter, r *http.Request) shttp.Response {
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

	brandModels, err := h.service.GetBrandModels(r.Context(), limit, page, search)
	if err != nil {
		h.logger.Error("unable to get brand models", err, err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(brandModels)
}

// v1UpdateBrandModel
// @Summary Update an existing brand model
// @Description Updates brand model details by ID
// @Tags Brand Model
// @Accept json
// @Produce json
// @Param brand body dtos.V1BrandModelDTO true "Brand Model data with ID"
// @Success 200 {object} map[string]int64 "Returns updated brand model ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/update-model [put]
func (h *BrandHandler) v1UpdateModel(w http.ResponseWriter, r *http.Request) shttp.Response {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(errBody.Error())
	}
	defer r.Body.Close()

	var modelDTO dtos.V1BrandModelDTO
	errData := json.Unmarshal(body, &modelDTO)
	if errData != nil {
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(errData.Error())
	}

	id, err := h.service.UpdateBrandModel(r.Context(), modelDTO)
	if err != nil {
		h.logger.Error("unable to update brand model", err, err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(map[string]interface{}{
		"id": id,
	})
}

// v1DeleteBrandModel
// @Summary Delete a brand model
// @Description Deletes a brand model by ID
// @Tags Brand Model
// @Accept json
// @Produce json
// @Param id query int true "Brand Model ID to delete"
// @Success 200 {object} string "brand model deleted successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Brand Model not found"
// @Failure 500 {object} string "Internal server error"
// @Router /brand/delete-model [delete]
func (h *BrandHandler) v1DeleteModel(w http.ResponseWriter, r *http.Request) shttp.Response {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return shttp.BadRequest.SetData("missing brand model ID")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("invalid brand model ID", err)
		return shttp.BadRequest.SetData("invalid brand model ID")
	}

	err = h.service.DeleteBrand(r.Context(), id)
	if err != nil {
		h.logger.Error("unable to delete brand model", err)
		return shttp.InternalServerError.SetData(err.Error())
	}

	return shttp.Success.SetData("brand model deleted successfully")
}
