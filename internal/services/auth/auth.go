package auth

import (
	"context"
	"grubzo/internal/config"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"
	"grubzo/internal/services/rbac"
	"grubzo/internal/services/user"
	"grubzo/internal/utils"

	"go.uber.org/zap"
)

type AuthService interface {
	BasicUserLogin(email, password string, tenantID uint) (uint, error)
	BasicEmployeeLogin(email, password string, tenantID uint) (uint, error)
	GetMeInfo(c context.Context, userType string, userID, tenantID uint, LocationID uint) (*dto.MeResponse, error)
}

type authServiceImpl struct {
	repo          *repository.Repository
	config        *config.Config
	userService   user.UserService
	rbacService   rbac.RBAC
	logger        *zap.Logger
}

func InitAuthService(
	repo *repository.Repository,
	config *config.Config,
	userService user.UserService,
	rbacService rbac.RBAC,
	logger *zap.Logger,
) (AuthService, error) {
	return &authServiceImpl{
		repo:          repo,
		config:        config,
		userService:   userService,
		rbacService:   rbacService,
		logger:        logger.Named("auth_service"),
	}, nil
}

func (a *authServiceImpl) BasicUserLogin(email, password string, tenantID uint) (uint, error) {
	u, err := a.repo.FindUser(query.NewUserQuery(tenantID).WithEmail(email))
	if err != nil {
		a.logger.Warn("user not found", zap.String("email", email), zap.Uint("tenantID", tenantID))
		return 0, ext.Error("invalid username or password")
	}
	if !a.repo.CheckUserPassword(u, password) {
		a.logger.Warn("invalid password", zap.String("email", email))
		return 0, ext.Error("invalid username or password")
	}
	return u.ID, nil
}

func (a *authServiceImpl) BasicEmployeeLogin(email, password string, tenantID uint) (uint, error) {
	u, err := a.repo.FindTenantUser(query.NewTenantUserQuery(tenantID).WithEmail(email))
	if err != nil {
		a.logger.Warn("user not found", zap.String("email", email), zap.Uint("tenantID", tenantID))
		return 0, ext.Error("invalid username or password")
	}
	if !a.repo.CheckTenantUserPassword(u, password) {
		a.logger.Warn("invalid password", zap.String("email", email))
		return 0, ext.Error("invalid username or password")
	}
	return u.ID, nil
}

func (a *authServiceImpl) GetMeInfo(c context.Context, userType string, userID, tenantID uint, locationID uint) (*dto.MeResponse, error) {
	me := &dto.MeResponse{
		Type: userType,
		ID:   userID,
	}
	q := query.NewTenantLocationQuery(tenantID)
	if locationID == 0 {
		q.WithPrimary(true)
	}else{
		q.WithID(locationID)
	}
	location, err := a.repo.FindTenantLocation(q)
	if err != nil{
		a.logger.Error("error fetching tenant locations", zap.Error(err))
	}
	locationInfo := dto.TenantLocation{}
	utils.Map(&locationInfo, location)
	me.Location = locationInfo

	if userType == "user" {
		user, err := a.repo.FindUser(query.NewUserQuery(tenantID).WithID(userID))
		if err != nil {
			return nil, err
		}
		me.Email = user.Email
		me.Name = user.Name
		return me, nil
	}
	if userType == "employee" {
		user, err := a.repo.FindTenantUser(query.NewTenantUserQuery(tenantID).WithID(userID))
		if err != nil {
			return nil, err
		}
		permissions := a.rbacService.GetGrantedPermissionsForRoles(user.TenantID, user.Roles)
		me.Email = user.Email
		me.Name = user.Name
		me.Permisssions = permissions
		return me, nil
	}
	return nil, nil
}
