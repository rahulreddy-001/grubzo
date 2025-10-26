package tenant

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

func (ts *tenantServiceImpl) CreateTenantUser(tUserArgs *dto.CreateTenantUser) (*dto.CreateTenantUserResponse, error) {
	user, err := ts.repository.CreateTenantUser(tUserArgs)
	if err != nil {
		return nil, err
	}
	userInfo := dto.TenantUserInfo{}
	if copier.Copy(&userInfo, user) != nil {
		ts.logger.Error("[copier.Copy] failed to copy user to userInfo", zap.Any("user", user), zap.Any("userInfo", userInfo))
	}
	response := &dto.CreateTenantUserResponse{
		Message: "User created successfully",
		User:    userInfo,
	}
	return response, nil

}

func (ts *tenantServiceImpl) UpdateTenantUser(args *dto.UpdateTenantUser) (*dto.UpdateTenantUserResponse, error) {
	user, err := ts.repository.UpdateTenantUser(args)
	if err != nil {
		return nil, err
	}
	userInfo := dto.TenantUserInfo{}
	if copier.Copy(&userInfo, user) != nil {
		ts.logger.Error("[copier.Copy] failed to copy user to userInfo", zap.Any("user", user), zap.Any("userInfo", userInfo))
	}
	response := &dto.UpdateTenantUserResponse{
		Message: "User updated successfully",
		User:    userInfo,
	}
	return response, nil
}

func (ts *tenantServiceImpl) GetTenantUser(UserID uint, tenantID uint) (*dto.GetTenantUserResponse, error) {
	user, err := ts.repository.FindTenantUser(query.NewTenantUserQuery(tenantID).WithID(UserID))
	if err != nil {
		return nil, err
	}
	userInfo := dto.TenantUserInfo{}
	if copier.Copy(&userInfo, user) != nil {
		ts.logger.Error("[copier.Copy] failed to copy user to userInfo", zap.Any("user", user), zap.Any("userInfo", userInfo))
	}
	response := &dto.GetTenantUserResponse{
		Message: "User fetched successfully",
		User:    userInfo,
	}
	return response, nil
}

func (ts *tenantServiceImpl) GetTenantUsers(tenantID uint) (*dto.GetTenantUsersResponse, error) {
	users, err := ts.repository.FindAllTenantUsers(query.NewTenantUserQuery(tenantID))
	if err != nil {
		return nil, err
	}
	usersInfo := make([]dto.TenantUserInfo, len(users))
	for i, user := range users {
		if copier.Copy(&usersInfo[i], user) != nil {
			ts.logger.Error("[copier.Copy] failed to copy user to usersInfo[i]", zap.Any("user", user), zap.Any("userInfo", usersInfo[i]))
		}

	}
	response := &dto.GetTenantUsersResponse{
		Message: "Users fetched successfully",
		Users:   usersInfo,
	}
	return response, nil
}
