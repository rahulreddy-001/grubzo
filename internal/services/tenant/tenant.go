package tenant

//go:generate go run ../../../cmd/injecttrace -file tenant.go -receiver tenantServiceImpl -service TenantService
import (
	"context"
	"grubzo/internal/config"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"
	"strings"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type TenantService interface {
	// Tenant
	CreateTenant(ctx context.Context, dto *dto.CreateTenant) (*dto.CreateTenantResponse, error)
	UpdateTenant(ctx context.Context, tenantID uint64, dto *dto.UpdateTenant) (*dto.UpdateTenantResponse, error)
	GetTenant(ctx context.Context, tenantID uint64) (*dto.GetTenantResponse, error)
	GetAllTenants(ctx context.Context) (*dto.GetAllTenantsResponse, error)
	GetTenants(ctx context.Context) (dto.GetSubDomainsResponse, error)

	// TenantLocation
	CreateTenantLocation(ctx context.Context, dto *dto.CreateTenantLocation) (*dto.CreateTenantLocationResponse, error)
	UpdateTenantLocation(ctx context.Context, dto *dto.UpdateTenantLocation) (*dto.UpdateTenantLocationResponse, error)
	GetTenantLocation(ctx context.Context, tenantLocId uint64, tenantID uint64) (*dto.TenantLocationResponse, error)
	GetTenantLocations(ctx context.Context, tenantID uint64) (*dto.TenantLocationsResponse, error)

	//TenantUser
	CreateTenantUser(ctx context.Context, dto *dto.CreateTenantUser) (*dto.CreateTenantUserResponse, error)
	UpdateTenantUser(ctx context.Context, dto *dto.UpdateTenantUser) (*dto.UpdateTenantUserResponse, error)
	GetTenantUser(ctx context.Context, UserID uint64, tenantID uint64) (*dto.GetTenantUserResponse, error)
	GetTenantUsers(ctx context.Context, tenantID uint64) (*dto.GetTenantUsersResponse, error)
	FetchTenantUsers(ctx context.Context, query *query.TenantUserQuery) (*dto.GetTenantUsersResponse, error)
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

func (ts *tenantServiceImpl) CreateTenant(ctx context.Context, args *dto.CreateTenant) (*dto.CreateTenantResponse, error) {
	ctx, span := otel.Tracer("TenantService").Start(ctx, "TenantService.CreateTenant")
	defer span.End()

	tenant, err := ts.repository.CreateTenant(ctx, args)
	if err != nil {
		return nil, err
	}
	tenantInfo := dto.TenantInfo{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Code:      tenant.Code,
		SubDomain: tenant.SubDomain,
	}
	response := &dto.CreateTenantResponse{
		Message: "Tenant created succssfully",
		Tenant:  tenantInfo,
	}
	return response, nil
}

func (ts *tenantServiceImpl) UpdateTenant(ctx context.Context, tenantID uint64, args *dto.UpdateTenant) (*dto.UpdateTenantResponse, error) {
	ctx, span := otel.Tracer("TenantService").Start(ctx, "TenantService.UpdateTenant")
	defer span.End()

	tenant, err := ts.repository.UpdateTenant(ctx, tenantID, args)
	if err != nil {
		return nil, err
	}
	tenantInfo := dto.TenantInfo{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Code:      tenant.Code,
		SubDomain: tenant.SubDomain,
	}
	response := &dto.UpdateTenantResponse{
		Message: "Tenant updated succssfully",
		Tenant:  tenantInfo,
	}
	return response, nil
}

func (ts *tenantServiceImpl) GetTenant(ctx context.Context, tenantID uint64) (*dto.GetTenantResponse, error) {
	ctx, span := otel.Tracer("TenantService").Start(ctx, "TenantService.GetTenant")
	defer span.End()

	tenant, err := ts.repository.GetTenant(ctx, query.NewTenantQuery().WithPreloads().WithID(tenantID))
	if err != nil {
		return nil, err
	}
	tenantInfo := dto.TenantInfo{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Code:      tenant.Code,
		SubDomain: tenant.SubDomain,
	}
	response := &dto.GetTenantResponse{
		Message: "Tenant fetched succssfully",
		Tenant:  tenantInfo,
	}
	return response, nil
}

func (ts *tenantServiceImpl) GetAllTenants(ctx context.Context) (*dto.GetAllTenantsResponse, error) {
	ctx, span := otel.Tracer("TenantService").Start(ctx, "TenantService.GetAllTenants")
	defer span.End()

	tenants, err := ts.repository.GetTenants(ctx, query.NewTenantQuery())
	if err != nil {
		return nil, err
	}
	tenantsInfo := []dto.TenantInfo{}
	for _, tenant := range tenants {
		tenantsInfo = append(tenantsInfo, dto.TenantInfo{
			ID:        tenant.ID,
			Name:      tenant.Name,
			Code:      tenant.Code,
			SubDomain: tenant.SubDomain,
		})
	}
	return &dto.GetAllTenantsResponse{
		Message: "Tenants fetched successfully",
		Tenants: tenantsInfo,
	}, nil
}

func (ts *tenantServiceImpl) GetTenants(ctx context.Context) (dto.GetSubDomainsResponse, error) {
	ctx, span := otel.Tracer("TenantService").Start(ctx, "TenantService.GetAllTenants")
	defer span.End()

	tenants, err := ts.repository.GetTenants(ctx, query.NewTenantQuery())
	if err != nil {
		return "", err
	}
	subDomains := []string{}
	for _, tenant := range tenants {
		subDomains = append(subDomains, tenant.SubDomain)
	}
	list := strings.Join(subDomains, ",")
	return dto.GetSubDomainsResponse(list), nil

}
