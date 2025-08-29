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

type SliderHandler struct {
	logger     *slog.Logger
	middleware *shttp.Middleware
	service    repository.SlidersService
}

func NewSliderHandler(logger *slog.Logger, middleware *shttp.Middleware, service repository.SlidersService) *SliderHandler {
	return &SliderHandler{
		logger:     logger,
		middleware: middleware,
		service:    service,
	}
}

func (h *SliderHandler) SliderRegisterRoutes(r chi.Router) {
	r.Method("POST", "/create-slider", h.middleware.Base(h.v1CreateSlider))
	r.Method("GET", "/get-sliders", h.middleware.Base(h.v1GetAllSliders))
	r.Method("PUT", "/update-slider", h.middleware.Base(h.v1UpdateSlider))
	r.Method("DELETE", "/delete-slider", h.middleware.Base(h.v1DeleteSlider))
}

// v1CreateSlider
// @Summary Create a new slider
// @Description Creates a new slider with the given name and image path
// @Tags Slider
// @Accept json
// @Produce json
// @Param slider body dtos.CreateSliderReq true "Slider data"
// @Success 200 {object} dtos.ID "Returns created slider ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /sliders/create-slider [post]
func (h *SliderHandler) v1CreateSlider(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var sliderDTO dtos.CreateSliderReq
	errData := json.Unmarshal(body, &sliderDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.CreateSlider(r.Context(), sliderDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to create slider", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Slider created Successfully"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1GetSliders
// @Summary Get all sliders
// @Description Get a paginated list of sliders with optional search
// @Tags Slider
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of sliders to return"
// @Param page query int false "Page number"
// @Param platform query string false "Platform string to filter sliders (web, mobile)"
// @Success 200 {object} dtos.SliderResult "List of sliders with pagination info successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /sliders/get-sliders [get]
func (h *SliderHandler) v1GetAllSliders(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")
	platform := r.URL.Query().Get("platform")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		limit = 10
	}
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page <= 0 {
		page = 1
	}

	sliders, err := h.service.GetAllSliders(r.Context(), limit, page, platform)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to get sliders", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "List of sliders with pagination info successfully"
	result.Data = sliders
	return shttp.Success.SetData(result)
}

// v1UpdateSlider handler
// @Summary Update an existing slider
// @Description Updates slider details by ID
// @Tags Slider
// @Accept json
// @Produce json
// @Param slider body dtos.UpdateSliderReq true "Slider data with ID"
// @Success 200 {object} dtos.ID "Returns updated slider ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /sliders/update-slider [put]
func (h *SliderHandler) v1UpdateSlider(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		result.Message = errBody.Error()
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(result)
	}
	defer r.Body.Close()

	var sliderDTO dtos.UpdateSliderReq
	errData := json.Unmarshal(body, &sliderDTO)
	if errData != nil {
		result.Message = errData.Error()
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(result)
	}

	id, err := h.service.UpdateSlider(r.Context(), sliderDTO)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to update slider", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Slider updated Successfully"
	result.Data = id
	return shttp.Success.SetData(result)
}

// v1DeleteSlider
// @Summary Delete a slider
// @Description Deletes a slider by ID
// @Tags Slider
// @Accept json
// @Produce json
// @Param id query int true "Slider ID to delete"
// @Success 200 {object} string "Slider deleted successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Brand not found"
// @Failure 500 {object} string "Internal server error"
// @Router /sliders/delete-slider [delete]
func (h *SliderHandler) v1DeleteSlider(w http.ResponseWriter, r *http.Request) shttp.Response {
	var result shttp.Result
	result.Status = false

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		result.Message = "missing slider ID"
		return shttp.BadRequest.SetData(result)
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("invalid slider ID", err)
		return shttp.BadRequest.SetData(result)
	}

	err = h.service.DeleteSlider(r.Context(), id)
	if err != nil {
		result.Message = err.Error()
		h.logger.Error("unable to delete slider", err)
		return shttp.InternalServerError.SetData(result)
	}

	result.Status = true
	result.Message = "Slider deleted Successfully"
	return shttp.Success.SetData(result)
}
