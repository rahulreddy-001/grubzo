package tenant

import (
	"context"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/utils"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

func (ts *tenantServiceImpl) CreateTenantLocation(ctx context.Context, tloc *dto.CreateTenantLocation) (*dto.CreateTenantLocationResponse, error) {
	locationEntity, err := ts.repository.CreateTenantLocation(ctx, tloc)
	if err != nil {
		return nil, err
	}
	location := dto.TenantLocation{}
	if utils.Map(&location, locationEntity) != nil {
		ts.logger.Error("[copier.Copy] failed to copy eLocation to tenantLocInfo", zap.Any("locationEntity", locationEntity), zap.Any("location", location), zap.Error(err))
	}
	response := &dto.CreateTenantLocationResponse{
		Message:  "Location created successfully",
		Location: location,
	}
	return response, nil
}

func (ts *tenantServiceImpl) UpdateTenantLocation(ctx context.Context, tloc *dto.UpdateTenantLocation) (*dto.UpdateTenantLocationResponse, error) {
	eLocation, err := ts.repository.UpdateTenantLocation(ctx, tloc)
	if err != nil {
		return nil, err
	}
	tenantLocInfo := dto.TenantLocation{}
	if copier.Copy(&tenantLocInfo, eLocation) != nil {
		ts.logger.Error("[copier.Copy] failed to copy eLocation to tenantLocInfo", zap.Any("eLocation", eLocation), zap.Any("tenantLocDTO", tenantLocInfo))
	}
	response := &dto.UpdateTenantLocationResponse{
		Message:  "Location updated successfully",
		Location: tenantLocInfo,
	}
	return response, nil
}

func (ts *tenantServiceImpl) GetTenantLocation(ctx context.Context, tenantLocId uint, tenantID uint) (*dto.TenantLocationResponse, error) {
	eLocation, err := ts.repository.FindTenantLocation(ctx, query.NewTenantLocationQuery(tenantID).WithID(tenantLocId))
	if err != nil {
		return nil, err
	}
	tenantLocInfo := dto.TenantLocation{}
	if copier.Copy(&tenantLocInfo, eLocation) != nil {
		ts.logger.Error("[copier.Copy] failed to copy eLocation to tenantLocInfo", zap.Any("eLocation", eLocation), zap.Any("tenantLocDTO", tenantLocInfo))
	}
	response := &dto.TenantLocationResponse{
		Message:  "Location fetched successfully",
		Location: tenantLocInfo,
	}
	return response, nil
}

func (ts *tenantServiceImpl) GetTenantLocations(ctx context.Context, tenantID uint) (*dto.TenantLocationsResponse, error) {
	eLocations, err := ts.repository.FindTenantLocations(ctx, query.NewTenantLocationQuery(tenantID))
	if err != nil {
		return nil, err
	}
	tenantLocInfos := make([]dto.TenantLocation, len(eLocations))
	for i, eLocation := range eLocations {
		if copier.Copy(&tenantLocInfos[i], eLocation) != nil {
			ts.logger.Error("[copier.Copy] failed to copy eLocation to tenantLocInfo", zap.Any("eLocation", eLocation), zap.Any("tenantLocDTO", tenantLocInfos[i]))
		}
	}
	response := &dto.TenantLocationsResponse{
		Message:   "Locations fetched successfully",
		Locations: tenantLocInfos,
	}
	return response, nil
}
