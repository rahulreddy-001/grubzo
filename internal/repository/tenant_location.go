package repository

import (
	"context"
	"errors"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"grubzo/internal/router/ext"

	"gorm.io/gorm"
)

type TenantLocationRepository interface {
	CreateTenantLocation(ctx context.Context, dto *dto.CreateTenantLocation) (*entity.TenantLocation, error)
	UpdateTenantLocation(ctx context.Context, dto *dto.UpdateTenantLocation) (*entity.TenantLocation, error)
	FindTenantLocation(ctx context.Context, query *query.TenantLocationQuery) (*entity.TenantLocation, error)
	FindTenantLocations(ctx context.Context, query *query.TenantLocationQuery) ([]*entity.TenantLocation, error)
}

func tenantLocationValidator(loc *entity.TenantLocation, db *gorm.DB) error {
	// Validate unique Code for tenant
	var existingByCode entity.TenantLocation
	if err := db.
		Where("tenant_id = ? AND code = ?", loc.TenantID, loc.Code).
		Not("id = ?", loc.ID).
		First(&existingByCode).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingByCode.ID != 0 {
		return ext.Error("A location with this code already exists.")
	}
	if loc.IsPrimary {
		var existingPrimary entity.TenantLocation
		if err := db.
			Where("tenant_id = ? AND is_primary = ?", loc.TenantID, true).
			Not("id = ?", loc.ID).
			First(&existingPrimary).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if existingPrimary.ID != 0 {
			return ext.Error("There cannot be more than one primary location.")
		}
	}

	return nil
}

func (r *Repository) CreateTenantLocation(ctx context.Context, dto *dto.CreateTenantLocation) (*entity.TenantLocation, error) {
	location := &entity.TenantLocation{
		TenantID:  dto.TenantID,
		Code:      dto.Code,
		Address:   dto.Address,
		City:      dto.City,
		State:     dto.State,
		Country:   dto.Country,
		ZipCode:   dto.ZipCode,
		IsPrimary: dto.IsPrimary,
	}

	db := r.dbWithContext(ctx)
	if err := tenantLocationValidator(location, db); err != nil {
		return nil, err
	}
	if err := db.Create(location).Error; err != nil {
		return nil, err
	}

	return location, nil
}

func (r *Repository) FindTenantLocation(ctx context.Context, q *query.TenantLocationQuery) (*entity.TenantLocation, error) {
	location := &entity.TenantLocation{}
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.TenantLocation{})

	sess = sess.Where("tenant_id = ?", q.TenantID)
	if q.ID != nil {
		sess = sess.Where("id = ?", *q.ID)
	}
	if q.Code != nil {
		sess = sess.Where("code = ?", *q.Code)
	}
	if q.IsPrimary != nil {
		sess = sess.Where("is_primary = ?", *q.IsPrimary)
	}
	if err := sess.First(&location).Error; err != nil {
		return nil, err
	}
	return location, nil
}

func (r *Repository) UpdateTenantLocation(ctx context.Context, dto *dto.UpdateTenantLocation) (*entity.TenantLocation, error) {
	location, err := r.FindTenantLocation(ctx, query.NewTenantLocationQuery(dto.TenantID).WithID(dto.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ext.Error("Location not found")
		}
	}

	if dto.Code != nil {
		location.Code = *dto.Code
	}
	if dto.Address != nil {
		location.Address = *dto.Address
	}
	if dto.City != nil {
		location.City = *dto.City
	}
	if dto.State != nil {
		location.State = *dto.State
	}
	if dto.Country != nil {
		location.Country = *dto.Country
	}
	if dto.ZipCode != nil {
		location.ZipCode = *dto.ZipCode
	}
	if dto.IsPrimary != nil {
		location.IsPrimary = *dto.IsPrimary
	}

	db := r.dbWithContext(ctx)
	if err := tenantLocationValidator(location, db); err != nil {
		return nil, err
	}
	if err := db.Save(&location).Error; err != nil {
		return nil, err
	}
	return location, nil
}

func (r *Repository) FindTenantLocations(ctx context.Context, q *query.TenantLocationQuery) ([]*entity.TenantLocation, error) {
	var locations []*entity.TenantLocation
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.TenantLocation{})

	sess = sess.Where("tenant_id = ?", q.TenantID)
	if q.ID != nil {
		sess = sess.Where("id = ?", *q.ID)
	}
	if q.Code != nil {
		sess = sess.Where("code = ?", *q.Code)
	}
	if q.IsPrimary != nil {
		sess = sess.Where("is_primary = ?", *q.IsPrimary)
	}

	if err := sess.Find(&locations).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ext.Error("Location not found")
		}
		return nil, err
	}
	return locations, nil
}
