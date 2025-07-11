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

type SettingsHandler struct {
	logger     *slog.Logger
	middleware *shttp.Middleware
	service    repository.SettingsService
}

func NewSettingsHandler(logger *slog.Logger, middleware *shttp.Middleware, service repository.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		logger:     logger,
		middleware: middleware,
		service:    service,
	}
}

func (h *SettingsHandler) SettingsRegisterRoutes(r chi.Router) {
	r.Method("POST", "/create-role", h.middleware.Base(h.v1CreateRole))
	r.Method("GET", "/get-roles", h.middleware.Base(h.v1GetAllRoles))
	r.Method("PUT", "/update-role", h.middleware.Base(h.v1UpdateRole))
	r.Method("DELETE", "/delete-role", h.middleware.Base(h.v1DeleteRole))

	//Users
	r.Method("POST", "/create-user", h.middleware.Base(h.v1CreateUser))
	r.Method("GET", "/get-users", h.middleware.Base(h.v1GetAllUsers))
	r.Method("PUT", "/update-user", h.middleware.Base(h.v1UpdateUser))
	r.Method("DELETE", "/delete-user", h.middleware.Base(h.v1DeleteUser))
}

// v1CreateRole
// @Summary Create a new role
// @Description Creates a new role with the given name and role
// @Tags Role
// @Accept json
// @Produce json
// @Param Role body dtos.Role true "Role data (the 'role' field accepts any JSON object)"
// @Success 200 {object} map[string]int64 "Returns created role ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /settings/create-role [post]
func (h *SettingsHandler) v1CreateRole(w http.ResponseWriter, r *http.Request) shttp.Response {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(errBody.Error())
	}
	defer r.Body.Close()

	var roleDTO dtos.Role
	errData := json.Unmarshal(body, &roleDTO)
	if errData != nil {
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(errData.Error())
	}

	id, err := h.service.CreateRole(r.Context(), roleDTO)
	if err != nil {
		h.logger.Error("unable to create role", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(map[string]interface{}{
		"id": id,
	})
}

// v1GetAllRoles
// @Summary Get all roles
// @Description Get a paginated list of roles with optional search
// @Tags Role
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of roles to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter roles by name"
// @Success 200 {object} dtos.RoleResult "List of roles with pagination info"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /settings/get-roles [get]
func (h *SettingsHandler) v1GetAllRoles(w http.ResponseWriter, r *http.Request) shttp.Response {
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

	roles, err := h.service.GetAllRoles(r.Context(), limit, page, search)
	if err != nil {
		h.logger.Error("unable to get roles", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(roles)
}

// v1UpdateRole handler
// @Summary Update an existing role
// @Description Updates role details by ID
// @Tags Role
// @Accept json
// @Produce json
// @Param Role body dtos.Role true "Role data with ID"
// @Success 200 {object} map[string]int64 "Returns updated role ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /settings/update-role [put]
func (h *SettingsHandler) v1UpdateRole(w http.ResponseWriter, r *http.Request) shttp.Response {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(errBody.Error())
	}
	defer r.Body.Close()

	var roleDTO dtos.Role
	errData := json.Unmarshal(body, &roleDTO)
	if errData != nil {
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(errData.Error())
	}

	id, err := h.service.UpdateRole(r.Context(), roleDTO)
	if err != nil {
		h.logger.Error("unable to update role", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(map[string]interface{}{
		"id": id,
	})
}

// v1DeleteRole
// @Summary Delete a role
// @Description Delete a role by ID
// @Tags Role
// @Accept json
// @Produce json
// @Param id query int true "Role ID to delete"
// @Success 200 {object} string "Role deleted successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Brand not found"
// @Failure 500 {object} string "Internal server error"
// @Router /settings/delete-role [delete]
func (h *SettingsHandler) v1DeleteRole(w http.ResponseWriter, r *http.Request) shttp.Response {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return shttp.BadRequest.SetData("missing role ID")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("invalid role ID", err)
		return shttp.BadRequest.SetData("invalid role ID")
	}

	err = h.service.DeleteRole(r.Context(), id)
	if err != nil {
		h.logger.Error("unable to delete role", err)
		return shttp.InternalServerError.SetData(err.Error())
	}

	return shttp.Success.SetData("role deleted successfully")
}

// v1CreateUser
// @Summary Create a new user
// @Description Creates a new user
// @Tags User
// @Accept json
// @Produce json
// @Param Role body dtos.User true "User data"
// @Success 200 {object} map[string]int64 "Returns created user ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /settings/create-user [post]
func (h *SettingsHandler) v1CreateUser(w http.ResponseWriter, r *http.Request) shttp.Response {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(errBody.Error())
	}
	defer r.Body.Close()

	var userDTO dtos.User
	errData := json.Unmarshal(body, &userDTO)
	if errData != nil {
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(errData.Error())
	}

	id, err := h.service.CreateUser(r.Context(), userDTO)
	if err != nil {
		h.logger.Error("unable to create user", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(map[string]interface{}{
		"id": id,
	})
}

// v1GetAllUsers
// @Summary Get all users
// @Description Get a paginated list of users with optional search
// @Tags User
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of users to return"
// @Param page query int false "Page number"
// @Param search query string false "Search string to filter users by name or login or roles by name"
// @Success 200 {object} dtos.UserResult "List of users with pagination info"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /settings/get-users [get]
func (h *SettingsHandler) v1GetAllUsers(w http.ResponseWriter, r *http.Request) shttp.Response {
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

	users, err := h.service.GetAllUsers(r.Context(), limit, page, search)
	if err != nil {
		h.logger.Error("unable to get users", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(users)
}

// v1UpdateUser handler
// @Summary Update an existing user
// @Description Updates user details by ID
// @Tags User
// @Accept json
// @Produce json
// @Param Role body dtos.User true "User data with ID"
// @Success 200 {object} map[string]int64 "Returns updated user ID"
// @Failure 400 {object} string "Bad request"
// @Failure 422 {object} string "Unprocessable entity"
// @Failure 500 {object} string "Internal server error"
// @Router /settings/update-user [put]
func (h *SettingsHandler) v1UpdateUser(w http.ResponseWriter, r *http.Request) shttp.Response {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		h.logger.Error("unable to read request body", errBody)
		return shttp.BadRequest.SetData(errBody.Error())
	}
	defer r.Body.Close()

	var userDTO dtos.User
	errData := json.Unmarshal(body, &userDTO)
	if errData != nil {
		h.logger.Error("unable to unmarshal request body", errData)
		return shttp.UnprocessableEntity.SetData(errData.Error())
	}

	id, err := h.service.UpdateUser(r.Context(), userDTO)
	if err != nil {
		h.logger.Error("unable to update user", err)
		return shttp.InternalServerError.SetData(err.Error())
	}
	return shttp.Success.SetData(map[string]interface{}{
		"id": id,
	})
}

// v1DeleteUser
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id query int true "User ID to delete"
// @Success 200 {object} string "User deleted successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Brand not found"
// @Failure 500 {object} string "Internal server error"
// @Router /settings/delete-user [delete]
func (h *SettingsHandler) v1DeleteUser(w http.ResponseWriter, r *http.Request) shttp.Response {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return shttp.BadRequest.SetData("missing user ID")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("invalid role ID", err)
		return shttp.BadRequest.SetData("invalid user ID")
	}

	err = h.service.DeleteRole(r.Context(), id)
	if err != nil {
		h.logger.Error("unable to delete user", err)
		return shttp.InternalServerError.SetData(err.Error())
	}

	return shttp.Success.SetData("user deleted successfully")
}
