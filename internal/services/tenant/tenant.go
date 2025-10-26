package tenant

import (
	"grubzo/internal/config"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"

	"go.uber.org/zap"
)

type TenantService interface {
	// Tenant
	CreateTenant(dto *dto.CreateTenant) (*dto.CreateTenantResponse, error)
	UpdateTenant(dto *dto.UpdateTenant) (*dto.UpdateTenantResponse, error)
	GetTenant(tenantID uint) (*dto.GetTenantResponse, error)
	GetAllTenants() (*dto.GetAllTenantsResponse, error)

	// TenantLocation
	CreateTenantLocation(dto *dto.CreateTenantLocation) (*dto.CreateTenantLocationResponse, error)
	UpdateTenantLocation(dto *dto.UpdateTenantLocation) (*dto.UpdateTenantLocationResponse, error)
	GetTenantLocation(tenantLocId uint, tenantID uint) (*dto.TenantLocationResponse, error)
	GetTenantLocations(tenantID uint) (*dto.TenantLocationsResponse, error)

	//TenantUser
	CreateTenantUser(dto *dto.CreateTenantUser) (*dto.CreateTenantUserResponse, error)
	UpdateTenantUser(dto *dto.UpdateTenantUser) (*dto.UpdateTenantUserResponse, error)
	GetTenantUser(UserID uint, tenantID uint) (*dto.GetTenantUserResponse, error)
	GetTenantUsers(tenantID uint) (*dto.GetTenantUsersResponse, error)
}

type tenantServiceImpl struct {
	repository *repository.Repository
	config     *config.Config
	logger     *zap.Logger
}

func InitTenantService(repository *repository.Repository, config *config.Config, logger *zap.Logger) (*tenantServiceImpl, error) {
	return &tenantServiceImpl{
		repository: repository,
		config:     config,
		logger:     logger.Named("tenant_service"),
	}, nil
}

func (ts *tenantServiceImpl) CreateTenant(args *dto.CreateTenant) (*dto.CreateTenantResponse, error) {
	tenant, err := ts.repository.CreateTenant(args)
	if err != nil {
		return nil, err
	}
	tenantInfo := dto.TenantInfo{
		ID:   tenant.ID,
		Name: tenant.Name,
	}
	response := &dto.CreateTenantResponse{
		Message: "Tenant created succssfully",
		Tenant:  tenantInfo,
	}
	return response, nil
}

func (ts *tenantServiceImpl) UpdateTenant(args *dto.UpdateTenant) (*dto.UpdateTenantResponse, error) {
	tenant, err := ts.repository.UpdateTenant(args)
	if err != nil {
		return nil, err
	}
	tenantInfo := dto.TenantInfo{
		ID:   tenant.ID,
		Name: tenant.Name,
	}
	response := &dto.UpdateTenantResponse{
		Message: "Tenant updated succssfully",
		Tenant:  tenantInfo,
	}
	return response, nil
}

func (ts *tenantServiceImpl) GetTenant(tenantID uint) (*dto.GetTenantResponse, error) {
	tenant, err := ts.repository.GetTenant(query.NewTenantQuery().WithPreloads().WithID(tenantID))
	if err != nil {
		return nil, err
	}
	tenantInfo := dto.TenantInfo{
		ID:   tenant.ID,
		Name: tenant.Name,
	}
	response := &dto.GetTenantResponse{
		Message: "Tenant fetched succssfully",
		Tenant:  tenantInfo,
	}
	return response, nil
}

func (ts *tenantServiceImpl) GetAllTenants() (*dto.GetAllTenantsResponse, error) {
	tenants, err := ts.repository.GetTenants(query.NewTenantQuery())
	if err != nil {
		return nil, err
	}
	tenantsInfo := []dto.TenantInfo{}
	for _, tenant := range tenants {
		tenantsInfo = append(tenantsInfo, dto.TenantInfo{
			ID:   tenant.ID,
			Name: tenant.Name,
		})
	}
	return &dto.GetAllTenantsResponse{
		Message: "Tenants fetched successfully",
		Tenants: tenantsInfo,
	}, nil
}
