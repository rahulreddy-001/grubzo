package user

import (
	"grubzo/internal/config"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(dto *dto.CreateUser) (*dto.CreateUserResponse, error)
	UpdateUser(dto *dto.UpdateUser) (*dto.UpdateUserResponse, error)
	GetUser(UserID uint, tenantID uint) (*dto.GetUserResponse, error)
	GetUsers(tenantID uint) (*dto.GetUsersResponse, error)
}

type userServiceImpl struct {
	repository *repository.Repository
	config     *config.Config
	logger     *zap.Logger
}

func InitUserService(repository *repository.Repository, config *config.Config, logger *zap.Logger) (*userServiceImpl, error) {
	return &userServiceImpl{
		repository: repository,
		config:     config,
		logger:     logger.Named("user_service"),
	}, nil
}

func (us *userServiceImpl) CreateUser(args *dto.CreateUser) (*dto.CreateUserResponse, error) {
	user, err := us.repository.CreateUser(args)
	if err != nil {
		return nil, err
	}
	userInfo := &dto.UserInfo{}
	if copier.Copy(userInfo, user) != nil {
		us.logger.Error("[copier.Copy] failed to copy user to userInfo", zap.Any("user", user), zap.Any("userInfo", userInfo))
	}
	response := &dto.CreateUserResponse{
		Message: "User created successfully",
		User:    *userInfo,
	}
	return response, nil
}

func (us *userServiceImpl) UpdateUser(args *dto.UpdateUser) (*dto.UpdateUserResponse, error) {
	user, err := us.repository.UpdateUser(args)
	if err != nil {
		return nil, err
	}
	userInfo := dto.UserInfo{}
	if copier.Copy(&userInfo, user) != nil {
		us.logger.Error("[copier.Copy] failed to copy user to userInfo", zap.Any("user", user), zap.Any("userInfo", userInfo))
	}
	response := &dto.UpdateUserResponse{
		Message: "User updated successfully",
		User:    userInfo,
	}
	return response, nil
}

func (us *userServiceImpl) GetUser(UserID uint, tenantID uint) (*dto.GetUserResponse, error) {
	user, err := us.repository.FindUser(query.NewUserQuery(tenantID).WithID(UserID))
	if err != nil {
		return nil, err
	}
	userInfo := dto.UserInfo{}
	if copier.Copy(&userInfo, user) != nil {
		us.logger.Error("[copier.Copy] failed to copy user to userInfo", zap.Any("user", user), zap.Any("userInfo", userInfo))
	}
	response := &dto.GetUserResponse{
		Message: "User fetched successfully",
		User:    userInfo,
	}
	return response, nil
}

func (us *userServiceImpl) GetUsers(tenantID uint) (*dto.GetUsersResponse, error) {
	users, err := us.repository.FindAllUsers(query.NewUserQuery(tenantID))
	if err != nil {
		return nil, err
	}
	usersInfo := make([]dto.UserInfo, len(users))
	for i, user := range users {
		if copier.Copy(usersInfo[i], user) != nil {
			us.logger.Error("[copier.Copy] failed to copy user to usersInfo[i]", zap.Any("user", user), zap.Any("userInfo", usersInfo[i]))
		}

	}
	response := &dto.GetUsersResponse{
		Message: "Users fetched successfully",
		Users:   usersInfo,
	}
	return response, nil
}
