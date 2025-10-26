package repository

import (
	"errors"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"grubzo/internal/utils/ce"

	"gorm.io/gorm"
)

type TenantLocationRepository interface {
	CreateTenantLocation(dto *dto.CreateTenantLocation) (*entity.TenantLocation, error)
	UpdateTenantLocation(dto *dto.UpdateTenantLocation) (*entity.TenantLocation, error)
	FindTenantLocation(query *query.TenantLocationQuery) (*entity.TenantLocation, error)
	FindTenantLocations(query *query.TenantLocationQuery) ([]*entity.TenantLocation, error)
}

func tenantLocationValidator(loc *entity.TenantLocation, db *gorm.DB) error {
	// Check for unique code & one primary location
	sess := db.Session(&gorm.Session{}).Model(&entity.TenantLocation{})
	var existing []*entity.TenantLocation
	sess.
		Where("tenant_id = ?", loc.TenantID).
		Where("code = ?", loc.Code)
	if loc.ID != 0 {
		sess.Not("id = ?", loc.ID)
	}
	if err := sess.Find(&existing).Error; err != nil {
		return err
	}
	for _, location := range existing {
		if location.IsPrimary && loc.IsPrimary {
			return ce.New("There cannot be more than one primary location")
		}
		if location.Code == loc.Code {
			return ce.New("A location with this code already exists")
		}
	}
	return nil
}

func (r *Repository) CreateTenantLocation(dto *dto.CreateTenantLocation) (*entity.TenantLocation, error) {
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

	if err := tenantLocationValidator(location, r.db); err != nil {
		return nil, err
	}
	if err := r.db.Create(location).Error; err != nil {
		return nil, err
	}

	return location, nil
}

func (r *Repository) FindTenantLocation(q *query.TenantLocationQuery) (*entity.TenantLocation, error) {
	location := &entity.TenantLocation{}
	sess := r.db.Session(&gorm.Session{}).Model(&entity.TenantLocation{})

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

func (r *Repository) UpdateTenantLocation(dto *dto.UpdateTenantLocation) (*entity.TenantLocation, error) {
	location, err := r.FindTenantLocation(query.NewTenantLocationQuery(dto.TenantID).WithID(dto.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ce.New("Location not found")
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

	if err := tenantLocationValidator(location, r.db); err != nil {
		return nil, err
	}
	if err := r.db.Save(&location).Error; err != nil {
		return nil, err
	}
	return location, nil
}

func (r *Repository) FindTenantLocations(q *query.TenantLocationQuery) ([]*entity.TenantLocation, error) {
	var locations []*entity.TenantLocation
	sess := r.db.Session(&gorm.Session{}).Model(&entity.TenantLocation{})

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
			return nil, ce.New("Location not found")
		}
		return nil, err
	}
	return locations, nil
}
