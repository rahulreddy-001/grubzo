package user

//go:generate go run ../../../cmd/injecttrace -file user.go -receiver userServiceImpl -service UserService

import (
	"context"
	"grubzo/internal/config"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"

	"github.com/jinzhu/copier"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(ctx context.Context, dto *dto.CreateUser) (*dto.CreateUserResponse, error)
	UpdateUser(ctx context.Context, dto *dto.UpdateUser) (*dto.UpdateUserResponse, error)
	GetUser(ctx context.Context, UserID uint64, tenantID uint64) (*dto.GetUserResponse, error)
	GetUsers(ctx context.Context, tenantID uint64) (*dto.GetUsersResponse, error)
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

func (us *userServiceImpl) CreateUser(ctx context.Context, args *dto.CreateUser) (*dto.CreateUserResponse, error) {
	ctx, span := otel.Tracer("UserService").Start(ctx, "UserService.CreateUser")
	defer span.End()

	user, err := us.repository.CreateUser(ctx, args)
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

func (us *userServiceImpl) UpdateUser(ctx context.Context, args *dto.UpdateUser) (*dto.UpdateUserResponse, error) {
	ctx, span := otel.Tracer("UserService").Start(ctx, "UserService.UpdateUser")
	defer span.End()

	user, err := us.repository.UpdateUser(ctx, args)
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

func (us *userServiceImpl) GetUser(ctx context.Context, UserID uint64, tenantID uint64) (*dto.GetUserResponse, error) {
	ctx, span := otel.Tracer("UserService").Start(ctx, "UserService.GetUser")
	defer span.End()

	user, err := us.repository.FindUser(ctx, query.NewUserQuery(tenantID).WithID(UserID))
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

func (us *userServiceImpl) GetUsers(ctx context.Context, tenantID uint64) (*dto.GetUsersResponse, error) {
	ctx, span := otel.Tracer("UserService").Start(ctx, "UserService.GetUsers")
	defer span.End()

	users, err := us.repository.FindAllUsers(ctx, query.NewUserQuery(tenantID))
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
