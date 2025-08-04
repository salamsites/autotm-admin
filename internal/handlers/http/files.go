package http

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"autotm-admin/internal/services/repository"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	shttp "github.com/salamsites/package-http"
	slog "github.com/salamsites/package-log"
	"io"
	"net/http"
)

type FilesHandler struct {
	logger     *slog.Logger
	middleware *shttp.Middleware
	service    repository.FilesService
}

func NewFilesHandler(logger *slog.Logger, middleware *shttp.Middleware, service repository.FilesService) *FilesHandler {
	return &FilesHandler{
		logger:     logger,
		middleware: middleware,
		service:    service,
	}
}

func (h *FilesHandler) FilesRegisterRoutes(r chi.Router) {
	r.Method("POST", "/upload-image", h.middleware.Base(h.v1UploadImage))
	r.Method("POST", "/delete-image", h.middleware.Base(h.v1DeleteImage))
}

// v1UploadImage
// @Summary Upload an image(s)
// @Description Uploads one or multiple image files and returns the file path(s)
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file(s) to upload (single or multiple)"
// @Success 200 {object} dtos.ImagePath "Returns the uploaded image path(s)"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /files/upload-image [post]
func (h *FilesHandler) v1UploadImage(w http.ResponseWriter, r *http.Request) shttp.Response {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.logger.Error("failed to parse multipart form", err)
		return shttp.BadRequest.SetData("Invalid form data")
	}
	formFiles := r.MultipartForm.File["image"]

	if len(formFiles) == 0 {
		h.logger.Error("failed to parse multipart form")
		return shttp.BadRequest.SetData("No file(s) uploaded")
	}

	var dto dtos.ImagePath

	for _, fileHeader := range formFiles {
		file, err := fileHeader.Open()
		if err != nil {
			h.logger.Error("failed to open file for upload", err)
			continue
		}
		defer file.Close()

		imagePath, err := helpers.UploadImage(file, fileHeader)
		if err != nil {
			h.logger.Error("failed to upload image", err)
			continue
		}

		dto.Images = append(dto.Images, imagePath)
	}

	if len(dto.Images) == 1 {
		dto.ImagePath = dto.Images[0]
	}
	//file, fileHeader, err := r.FormFile("image")
	//if err != nil {
	//	h.logger.Error("unable to get uploaded file", err)
	//	return shttp.BadRequest.SetData(err.Error())
	//}
	//defer file.Close()
	//
	//imagePath, err := h.service.UploadImage(file, fileHeader)
	//if err != nil {
	//	h.logger.Error("unable to upload image error", err)
	//	return shttp.InternalServerError.SetData(err.Error())
	//}

	if len(dto.Images) == 0 {
		return shttp.InternalServerError.SetData("Failed to upload any images")
	}
	return shttp.Success.SetData(dto)
}

// v1DeleteImage
// @Summary Delete Image
// @Description Deletes Image
// @Tags Files
// @Accept json
// @Produce json
// @Param imagePath body dtos.ImagePath true "Image data"
// @Success 200 {object} map[string]int64 "Returns deleted Image"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /files/delete-image [post]
func (h *FilesHandler) v1DeleteImage(w http.ResponseWriter, r *http.Request) shttp.Response {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(errBody.Error())
	}
	defer r.Body.Close()

	var imagePath dtos.ImagePath
	errData := json.Unmarshal(body, &imagePath)
	if errData != nil {
		h.logger.Error("unable to unmarshal request body ", errData)
		return shttp.UnprocessableEntity.SetData(errData.Error())
	}
	if imagePath.ImagePath == "" && len(imagePath.Images) == 0 {
		h.logger.Error("no image paths provided to delete")
		return shttp.BadRequest.SetData("No image paths provided")
	}

	err := h.service.DeleteImage(r.Context(), imagePath)
	if err != nil {
		h.logger.Error("unable to delete image", err)
		return shttp.InternalServerError.SetData(err.Error())
	}

	return shttp.Success.SetData("image deleted successfully")
}
