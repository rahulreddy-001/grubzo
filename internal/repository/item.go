package repository

import (
	"errors"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"slices"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type ItemRepository interface {
	CreateItem(dto *dto.CreateMenuItem) (*entity.Item, error)
	UpdateItem(dto *dto.UpdateMenuItem) (*entity.Item, error)
	GetItem(filter *query.MenuItemQuery) (*entity.Item, error)
	GetItems(filter *query.MenuItemQuery) ([]*entity.Item, error)
}

func (r Repository) CreateItem(dto *dto.CreateMenuItem) (*entity.Item, error) {
	sess := r.db.Session(&gorm.Session{}).Model(&entity.Item{})
	item := &entity.Item{
		TenantID:          dto.TenantID,
		LocationID:        dto.LocationID,
		Name:              dto.Name,
		Description:       dto.Description,
		Price:             dto.Price,
		PriceUnit:         dto.PriceUnit,
		Category:          dto.Category,
		AvailableQuantity: dto.AvailableQuantity,
		Orderable:         dto.Orderable,
	}

	if err := sess.Create(item).Error; err != nil {
		return nil, err
	}

	if err := r.PopulateOwnerID(nil, item.ID, dto.Files, dto.TenantID); err != nil {
		if e := sess.Delete(item).Error; e != nil {
			return nil, errors.Join(e, err)
		}
		return nil, err
	}
	return r.GetItem(query.NewMenuItemQuery(item.TenantID).WithID(item.ID).WithPreload())
}

func (r *Repository) UpdateItem(dto *dto.UpdateMenuItem) (*entity.Item, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		item, err := r.GetItem(query.NewMenuItemQuery(dto.TenantID).WithID(dto.ID))
		if err != nil {
			return err
		}

		if dto.Name != nil {
			item.Name = *dto.Name
		}
		if dto.Description != nil {
			item.Description = *dto.Description
		}
		if dto.Price != nil {
			item.Price = *dto.Price
		}
		if dto.PriceUnit != nil {
			item.PriceUnit = *dto.PriceUnit
		}
		if dto.Category != nil {
			item.Category = *dto.Category
		}
		if dto.AvailableQuantity != nil {
			item.AvailableQuantity = *dto.AvailableQuantity
		}
		if dto.Orderable != nil {
			item.Orderable = *dto.Orderable
		}

		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		if len(dto.Files) > 0 {
			filesToDelete := []uuid.UUID{}
			for _, file := range item.Files {
				if !slices.Contains(dto.Files, file.ID) {
					filesToDelete = append(filesToDelete, file.ID)
				}
			}
			if err := r.DeleteFiles(tx, filesToDelete); err != nil {
				return err
			}

			if err := r.PopulateOwnerID(tx, item.ID, dto.Files, dto.TenantID); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetItem(query.NewMenuItemQuery(dto.TenantID).WithID(dto.ID).WithPreload())
}

func (r Repository) GetItem(filter *query.MenuItemQuery) (*entity.Item, error) {
	item := &entity.Item{}
	sess := r.db.Session(&gorm.Session{}).Model(&entity.Item{}).Where("tenant_id = ?", filter.TenantID)
	if filter.ID != nil {
		sess.Where("id = ?", filter.ID)
	}
	if filter.LocationID != nil {
		sess.Where("location_id = ?", filter.LocationID)
	}
	if filter.Orderable != nil {
		sess.Where("orderable = ?", filter.Orderable)
	}
	if filter.Preload {
		for _, preload := range item.GetPreloads() {
			sess.Preload(preload)
		}
	}
	if err := sess.Find(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (r Repository) GetItems(filter *query.MenuItemQuery) ([]*entity.Item, error) {
	items := []*entity.Item{}
	sess := r.db.Session(&gorm.Session{}).Model(&entity.Item{}).Where("tenant_id = ?", filter.TenantID)
	if filter.ID != nil {
		sess.Where("id = ?", filter.ID)
	}
	if filter.LocationID != nil {
		sess.Where("location_id = ?", filter.LocationID)
	}
	if filter.Orderable != nil {
		sess.Where("orderable = ?", filter.Orderable)
	}
	if filter.Preload {
		for _, preload := range (entity.Item{}).GetPreloads() {
			sess.Preload(preload)
		}
	}
	if err := sess.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
