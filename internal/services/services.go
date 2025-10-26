package services

import (
	"errors"
	"fmt"
	"grubzo/internal/config"
	"grubzo/internal/repository"
	"grubzo/internal/services/auth"
	"grubzo/internal/services/file"
	"grubzo/internal/services/store"
	"grubzo/internal/services/tenant"
	"grubzo/internal/services/user"
	"grubzo/internal/utils/storage"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Services struct {
	FileManager   file.Manager
	UserService   user.UserService
	TenantService tenant.TenantService
	StoreService  store.StoreService
	AuthService   auth.AuthService
}

func Setup(
	logger *zap.Logger,
	db *gorm.DB,
	repository *repository.Repository,
	fs storage.FileStorage,
	config *config.Config,
) (*Services, error) {
	services := &Services{}

	var errs []error

	fm, err := file.InitFileManager(repository, fs, logger)
	if err != nil {
		errs = append(errs, fmt.Errorf("InitFileManager: %w", err))
	} else {
		services.FileManager = fm
	}

	if us, err := user.InitUserService(repository, config, logger); err != nil {
		errs = append(errs, fmt.Errorf("InitUserService: %w", err))
	} else {
		services.UserService = us
	}

	if ts, err := tenant.InitTenantService(repository, config, logger); err != nil {
		errs = append(errs, fmt.Errorf("InitTenantService: %w", err))
	} else {
		services.TenantService = ts
	}

	if ss, err := store.Init(repository, config, fm, logger); err != nil {
		errs = append(errs, fmt.Errorf("InitStoreService: %w", err))
	} else {
		services.StoreService = ss
	}

	if as, err := auth.InitAuthService(repository, config, services.UserService, logger); err != nil {
		errs = append(errs, fmt.Errorf("InitAuthService: %w", err))
	} else {
		services.AuthService = as
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return services, nil
}
