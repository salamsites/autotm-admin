package http

import (
	"autotm-admin/internal/services/repository"
	"github.com/go-chi/chi/v5"
	shttp "github.com/salamsites/package-http"
	slog "github.com/salamsites/package-log"
	"net/http"
	"strconv"
)

type UsersHandler struct {
	logger     *slog.Logger
	middleware *shttp.Middleware
	service    repository.UserService
}

func NewUsersHandler(logger *slog.Logger, middleware *shttp.Middleware, service repository.UserService) *UsersHandler {
	return &UsersHandler{
		logger:     logger,
		middleware: middleware,
		service:    service,
	}
}

func (h *UsersHandler) UsersRegisterRoutes(r chi.Router) {
	r.Method("GET", "/get-users", h.middleware.Base(h.v1GetUsersFromUserService))
}

// v1GetUsersFromUserService
// @Summary Get all users from user service
// @Description Get a paginated list of users with optional search
// @Tags Users
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of users to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter users by name"
// @Success 200 {object} dtos.GetUsersResult "List of users with pagination info Successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /users/get-users [get]
func (h *UsersHandler) v1GetUsersFromUserService(w http.ResponseWriter, r *http.Request) shttp.Response {
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
	result.Message = "List of users with pagination info Successfully"
	result.Data = users
	return shttp.Success.SetData(result)
}
