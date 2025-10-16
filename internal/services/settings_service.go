package services

import (
	"autotm-admin/internal/configs"
	"autotm-admin/internal/dtos"
	"autotm-admin/internal/models"
	"autotm-admin/internal/repository/storage"
	"autotm-admin/utils"
	"context"

	slog "github.com/salamsites/package-log"
	"golang.org/x/crypto/bcrypt"
)

type SettingsService struct {
	logger *slog.Logger
	repo   storage.SettingsRepository
	cfg    *configs.Config
}

func NewSettingsService(logger *slog.Logger, repo storage.SettingsRepository, cfg *configs.Config) *SettingsService {
	return &SettingsService{
		logger: logger,
		repo:   repo,
		cfg:    cfg,
	}
}

func (s *SettingsService) CreateRole(ctx context.Context, role dtos.CreateRoleReq) (dtos.ID, error) {
	var id dtos.ID

	newRole := models.Role{
		Name: role.Name,
		Role: role.Role,
	}

	roleID, err := s.repo.CreateRole(ctx, newRole)
	if err != nil {
		s.logger.Errorf("create role service err: %v", err)
		return id, err
	}
	id.ID = roleID

	return id, nil
}

func (s *SettingsService) GetRoleByID(ctx context.Context, roleID int64) (dtos.Role, error) {
	role, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		s.logger.Errorf("get role by id service err: %v", err)
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

	roles, err := s.repo.GetAllRoles(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get all roles service err: %v", err)
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

	count, errCount := s.repo.GetRolesCount(ctx, search)
	if errCount != nil {
		s.logger.Errorf("get roles count service err: %v", err)
		return dtos.RoleResult{}, errCount
	}

	result := dtos.RoleResult{
		Roles: dtoRoles,
		Count: count,
	}

	return result, nil
}

func (s *SettingsService) UpdateRole(ctx context.Context, role dtos.UpdateRoleReq) (dtos.ID, error) {
	var id dtos.ID

	newRole := models.Role{
		ID:   role.ID,
		Name: role.Name,
		Role: role.Role,
	}

	roleID, err := s.repo.UpdateRole(ctx, newRole)
	if err != nil {
		s.logger.Errorf("update role service err: %v", err)
		return id, err
	}
	id.ID = roleID

	return id, nil
}

func (s *SettingsService) DeleteRole(ctx context.Context, id int64) error {
	err := s.repo.DeleteRole(ctx, id)
	if err != nil {
		s.logger.Errorf("delete role service err: %v", err)
		return err
	}
	return nil
}

// User
func (s *SettingsService) CreateUser(ctx context.Context, user dtos.CreateUserReq) (dtos.ID, error) {
	var id dtos.ID

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("hash err: %v", err)
		return id, err
	}

	newUser := models.User{
		Username: user.Username,
		Login:    user.Login,
		Password: string(hashedPassword),
		RoleID:   user.RoleID,
		Status:   user.Status,
	}

	userID, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		s.logger.Errorf("create user service err: %v", err)
		return id, err
	}
	id.ID = userID

	return id, nil
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
		s.logger.Errorf("create super admin service err: %v", err)
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

	users, err := s.repo.GetAllUsers(ctx, limit, offset, search)
	if err != nil {
		s.logger.Errorf("get all users service err: %v", err)
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
			Status:   b.Status,
		})
	}

	count, errCount := s.repo.GetUsersCount(ctx, search)
	if errCount != nil {
		s.logger.Errorf("get users count service err: %v", errCount)
		return dtos.UserResult{}, err
	}

	result := dtos.UserResult{
		Users: dtoUsers,
		Count: count,
	}

	return result, nil
}

func (s *SettingsService) UpdateUser(ctx context.Context, user dtos.UpdateUserReq) (dtos.ID, error) {
	var id dtos.ID

	var hashedPassword string

	if user.Password != "" {
		hp, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Errorf("hash err: %v", err)
			return id, err
		}
		hashedPassword = string(hp)
	}

	newUser := models.User{
		ID:       user.ID,
		Username: user.Username,
		Login:    user.Login,
		Password: hashedPassword,
		RoleID:   user.RoleID,
		Status:   user.Status,
	}

	userID, err := s.repo.UpdateUser(ctx, newUser)
	if err != nil {
		s.logger.Errorf("update user service  err: %v", err)
		return id, err
	}
	id.ID = userID
	return id, nil
}

func (s *SettingsService) DeleteUser(ctx context.Context, id int64) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		s.logger.Errorf("delete user service err: %v", err)
		return err
	}
	return nil
}

func (s *SettingsService) Login(ctx context.Context, login dtos.LoginReq) (string, error) {
	user, err := s.repo.GetUserByLogin(ctx, login.Login)
	if err != nil {
		s.logger.Errorf("get user login service err: %v", err)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		s.logger.Errorf("compare hash err: %v", err)
	}

	token, err := utils.TokenEncode(user.ID, user.RoleName, s.cfg.Auth.JwtRegistration)
	if err != nil {
		s.logger.Errorf("token encode err: %v", err)
		return "", err
	}

	return token, nil
}
