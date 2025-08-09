package services

import (
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/helpers"
	"autotm-admin/internal/models"
	"autotm-admin/internal/repository/storage"
	"context"
	slog "github.com/salamsites/package-log"
	"golang.org/x/crypto/bcrypt"
)

type SettingsService struct {
	logger *slog.Logger
	repo   storage.SettingsRepository
}

func NewSettingsService(logger *slog.Logger, repo storage.SettingsRepository) *SettingsService {
	return &SettingsService{
		logger: logger,
		repo:   repo,
	}
}

func (s *SettingsService) CreateRole(ctx context.Context, role dtos.CreateRoleReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(role); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	newRole := models.Role{
		Name: role.Name,
		Role: role.Role,
	}

	roleID, err := s.repo.CreateRole(ctx, newRole)
	if err != nil {
		s.logger.Errorf("create err: %v", err)
		return roleID, err
	}
	return roleID, nil
}

func (s *SettingsService) GetRoleByID(ctx context.Context, roleID int64) (dtos.Role, error) {
	role, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		s.logger.Errorf("get err: %v", err)
		return dtos.Role{}, err
	}

	result := dtos.Role{
		ID:   role.ID,
		Name: role.Name,
		Role: role.Role,
	}

	return result, nil
}

func (s *SettingsService) GetAllRoles(ctx context.Context, limit, page int64, search string) (dtos.RoleResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	roles, count, err := s.repo.GetAllRoles(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get all roles err: %v", err)
		return dtos.RoleResult{}, err
	}
	var dtoRoles []dtos.Role
	for _, b := range roles {
		dtoRoles = append(dtoRoles, dtos.Role{
			ID:   b.ID,
			Name: b.Name,
			Role: b.Role,
		})
	}

	result := dtos.RoleResult{
		Roles: dtoRoles,
		Count: count,
	}
	return result, nil
}

func (s *SettingsService) UpdateRole(ctx context.Context, role dtos.UpdateRoleReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(role); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}

	newRole := models.Role{
		ID:   role.ID,
		Name: role.Name,
		Role: role.Role,
	}

	roleID, err := s.repo.UpdateRole(ctx, newRole)
	if err != nil {
		s.logger.Errorf("update role err: %v", err)
		return roleID, err
	}
	return roleID, nil
}

func (s *SettingsService) DeleteRole(ctx context.Context, id int64) error {
	deleteID := models.ID{
		ID: id,
	}

	err := s.repo.DeleteRole(ctx, deleteID)
	if err != nil {
		s.logger.Errorf("delete role err: %v", err)
		return err
	}
	return nil
}

// User
func (s *SettingsService) CreateUser(ctx context.Context, user dtos.CreateUserReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(user); err != nil {
		s.logger.Errorf("validate user err: %v", err)
		return 0, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("hash err: %v", err)
		return 0, err
	}

	newUser := models.User{
		Username: user.Username,
		Login:    user.Login,
		Password: string(hashedPassword),
		RoleID:   user.RoleID,
	}

	userID, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		s.logger.Errorf("create user err: %v", err)
		return 0, err
	}
	return userID, nil
}

func (s *SettingsService) InitSuperAdmin(ctx context.Context) error {
	const (
		superAdminLogin    = "superadmin"
		superAdminPassword = "admin123"
		superUser          = "Super Admin"
	)
	existingUser, err := s.repo.GetUserByLogin(ctx, superAdminLogin)
	if err == nil && existingUser.ID != 0 {
		s.logger.Info("Super admin already exists.")
		return nil
	}

	superRole := models.Role{
		Name: superAdminLogin,
		Role: []byte(`{"permissions":"all"}`),
	}

	roleID, err := s.repo.CreateRole(ctx, superRole)
	if err != nil {
		s.logger.Errorf("create super role err: %v", err)
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(superAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("hash err: %v", err)
		return err
	}

	superAdmin := models.User{
		Username: superUser,
		Login:    superAdminLogin,
		Password: string(hashedPassword),
		RoleID:   roleID,
	}

	_, err = s.repo.CreateUser(ctx, superAdmin)
	if err != nil {
		s.logger.Errorf("create super admin err: %v", err)
		return err
	}

	s.logger.Info("Default super admin created successfully.")

	return nil
}

func (s *SettingsService) GetAllUsers(ctx context.Context, limit, page int64, search string) (dtos.UserResult, error) {
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
		offset = 0
	}

	users, count, err := s.repo.GetAllUsers(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get all users err: %v", err)
		return dtos.UserResult{}, err
	}
	var dtoUsers []dtos.User
	for _, b := range users {
		dtoUsers = append(dtoUsers, dtos.User{
			ID:       b.ID,
			Username: b.Username,
			Login:    b.Login,
			Password: b.Password,
			RoleID:   b.RoleID,
			RoleName: b.RoleName,
		})
	}

	result := dtos.UserResult{
		Users: dtoUsers,
		Count: count,
	}
	return result, nil
}

func (s *SettingsService) UpdateUser(ctx context.Context, user dtos.UpdateUserReq) (int64, error) {
	validate := helpers.GetValidator()
	if err := validate.Struct(user); err != nil {
		s.logger.Errorf("validate err: %v", err)
		return 0, err
	}
	var hashedPassword string

	if user.Password != "" {
		hp, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Errorf("hash err: %v", err)
			return 0, err
		}
		hashedPassword = string(hp)
	}

	newUser := models.User{
		ID:       user.ID,
		Username: user.Username,
		Login:    user.Login,
		Password: hashedPassword,
		RoleID:   user.RoleID,
	}

	userID, err := s.repo.UpdateUser(ctx, newUser)
	if err != nil {
		s.logger.Errorf("update user err: %v", err)
		return userID, err
	}
	return userID, nil
}

func (s *SettingsService) DeleteUser(ctx context.Context, id int64) error {
	deleteID := models.ID{
		ID: id,
	}

	err := s.repo.DeleteUser(ctx, deleteID)
	if err != nil {
		s.logger.Errorf("delete user err: %v", err)
		return err
	}
	return nil
}
