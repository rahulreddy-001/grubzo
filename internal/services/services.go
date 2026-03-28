package services

//go:generate go run ../../cmd/injecttrace -file services.go -receiver Services -service Services

import (
	"errors"
	"fmt"
	"grubzo/internal/config"
	"grubzo/internal/repository"
	"grubzo/internal/services/auth"
	"grubzo/internal/services/file"
	"grubzo/internal/services/order"
	"grubzo/internal/services/payment"
	"grubzo/internal/services/rbac"
	"grubzo/internal/services/store"
	"grubzo/internal/services/tenant"
	"grubzo/internal/services/user"
	"grubzo/internal/utils/storage"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Services struct {
	FileManager     file.Manager
	UserService     user.UserService
	TenantService   tenant.TenantService
	StoreService    store.StoreService
	AuthService     auth.AuthService
	CartService     order.CartService
	RBAC            *rbac.RBAC
	RazorpayService payment.RazorpayService
	WalletService   order.WalletService
	OrderService    order.OrderService
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

	if ac, err := rbac.New(repository); err != nil {
		errs = append(errs, fmt.Errorf("RABC: %w", err))
	} else {
		services.RBAC = ac
	}

	if as, err := auth.InitAuthService(repository, config, services.UserService, *services.RBAC, logger); err != nil {
		errs = append(errs, fmt.Errorf("InitAuthService: %w", err))
	} else {
		services.AuthService = as
	}

	if cs, err := order.InitCartService(repository, services.StoreService, config, logger); err != nil {
		errs = append(errs, fmt.Errorf("InitCartService: %w", err))
	} else {
		services.CartService = cs
	}

	if ps, err := payment.InitRazorpayService(repository, config, logger); err != nil {
		errs = append(errs, fmt.Errorf("InitPaymentService: %w", err))
	} else {
		services.RazorpayService = ps
	}

	if ws, err := order.InitWalletService(repository, services.RazorpayService, config, logger); err != nil {
		errs = append(errs, fmt.Errorf("InitWalletService: %w", err))
	} else {
		services.WalletService = ws
	}

	if os, err := order.InitOrderService(repository, services.WalletService, services.CartService, services.StoreService, config, logger); err != nil {
		errs = append(errs, fmt.Errorf("InitOrderService: %w", err))
	} else {
		services.OrderService = os
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return services, nil
}
