package http

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	sminio "github.com/salamsites/minio-pkg"
	"github.com/salamsites/minio-pkg/util"
	shttp "github.com/salamsites/package-http"
	slog "github.com/salamsites/package-log"
)

type FilesHandler struct {
	logger          *slog.Logger
	middleware      *shttp.Middleware
	minioFileClient sminio.ImageClient
}

func NewFilesHandler(logger *slog.Logger, middleware *shttp.Middleware, minioFileClient sminio.ImageClient) *FilesHandler {
	return &FilesHandler{
		logger:          logger,
		middleware:      middleware,
		minioFileClient: minioFileClient,
	}
}

func (h *FilesHandler) FilesRegisterRoutes(r chi.Router) {
	r.Method("POST", "/upload-image", h.middleware.Base(h.v1UploadImage))
	r.Method("DELETE", "/delete-image", h.middleware.Base(h.v1DeleteImage))
}

// v1UploadImage
// @Summary Upload an image
// @Description Uploads one image file and returns the file path
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file to upload (single)"
// @Success 200 {object} dtos.UploadImage "Returns the uploaded image path"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /files/upload-image [post]
func (h *FilesHandler) v1UploadImage(w http.ResponseWriter, r *http.Request) shttp.Response {
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		h.logger.Error("failed to parse multipart form", err)
		return shttp.BadRequest.SetData("Invalid form data")
	}

	fileHeaders := r.MultipartForm.File["image"]
	if len(fileHeaders) == 0 {
		h.logger.Error("no file uploaded")
		return shttp.BadRequest.SetData("No file uploaded")
	}

	fileHeader := fileHeaders[0]
	file, err := fileHeader.Open()
	if err != nil {
		h.logger.Error("failed to open uploaded file", err)
		return shttp.BadRequest.SetData("Failed to open file")
	}
	defer file.Close()

	uploadID := uuid.NewString()
	path := fmt.Sprintf("%s", uploadID)

	errUpload := h.minioFileClient.UploadImage(r.Context(), r, "image", path, helpers.FileSizes, util.FileBucket)
	if errUpload.StatusCode != 0 {
		h.logger.Error("failed to upload image", errUpload)
		return shttp.InternalServerError.SetData(errUpload.Message)
	}

	sizeStrings := make([]string, 0, len(helpers.FileSizes))
	for _, s := range helpers.FileSizes {
		if s.Height == 0 {
			sizeStrings = append(sizeStrings, fmt.Sprintf("%dx0", s.Width))
		} else {
			sizeStrings = append(sizeStrings, fmt.Sprintf("%dx%d", s.Width, s.Height))
		}
	}

	resp := dtos.UploadImage{
		UploadID: uploadID,
		Sizes:    sizeStrings,
	}

	return shttp.Success.SetData(resp)
}

// v1DeleteImage
// @Summary Delete Image
// @Description Deletes Image
// @Tags Files
// @Accept json
// @Produce json
// @Param upload_id query string true "Upload ID to delete"
// @Success 200 {object} string "Returns deleted Image"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /files/delete-image [delete]
func (h *FilesHandler) v1DeleteImage(w http.ResponseWriter, r *http.Request) shttp.Response {
	idStr := r.URL.Query().Get("upload_id")
	if idStr == "" {
		return shttp.BadRequest.SetData("id is required")
	}

	err := h.minioFileClient.RemoveImage(r.Context(), idStr, util.FileBucket)
	if err != nil {
		h.logger.Error("unable to delete image", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData("image deleted successfully")
}
