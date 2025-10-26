package tenant

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

func (ts *tenantServiceImpl) CreateTenantLocation(tloc *dto.CreateTenantLocation) (*dto.CreateTenantLocationResponse, error) {
	eLocation, err := ts.repository.CreateTenantLocation(tloc)
	if err != nil {
		return nil, err
	}
	tenantLocInfo := dto.TenantLocation{}
	if copier.Copy(&tenantLocInfo, eLocation) != nil {
		ts.logger.Error("[copier.Copy] failed to copy eLocation to tenantLocInfo", zap.Any("eLocation", eLocation), zap.Any("tenantLocDTO", tenantLocInfo))
	}
	response := &dto.CreateTenantLocationResponse{
		Message:  "Location created successfully",
		Location: tenantLocInfo,
	}
	return response, nil
}

func (ts *tenantServiceImpl) UpdateTenantLocation(tloc *dto.UpdateTenantLocation) (*dto.UpdateTenantLocationResponse, error) {
	eLocation, err := ts.repository.UpdateTenantLocation(tloc)
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

func (ts *tenantServiceImpl) GetTenantLocation(tenantLocId uint, tenantID uint) (*dto.TenantLocationResponse, error) {
	eLocation, err := ts.repository.FindTenantLocation(query.NewTenantLocationQuery(tenantID).WithID(tenantLocId))
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

func (ts *tenantServiceImpl) GetTenantLocations(tenantID uint) (*dto.TenantLocationsResponse, error) {
	eLocations, err := ts.repository.FindTenantLocations(query.NewTenantLocationQuery(tenantID))
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
