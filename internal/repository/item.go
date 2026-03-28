package repository

import (
	"context"
	"errors"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"slices"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type ItemRepository interface {
	CreateItem(ctx context.Context, dto *dto.CreateMenuItem) (*entity.Item, error)
	UpdateItem(ctx context.Context, dto *dto.UpdateMenuItem) (*entity.Item, error)
	GetItem(ctx context.Context, filter *query.MenuItemQuery) (*entity.Item, error)
	GetItems(ctx context.Context, filter *query.MenuItemQuery) ([]*entity.Item, error)
}

func (r Repository) CreateItem(ctx context.Context, dto *dto.CreateMenuItem) (*entity.Item, error) {
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.Item{})
	item := &entity.Item{
		TenantID:    dto.TenantID,
		LocationID:  dto.LocationID,
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Category:    dto.Category,
		FoodType:    dto.FoodType,
		ItemStatus:  dto.ItemStatus,
	}

	if err := sess.Create(item).Error; err != nil {
		return nil, err
	}

	if err := r.PopulateOwnerID(ctx, nil, item.ID, dto.Files, dto.TenantID); err != nil {
		if e := sess.Delete(item).Error; e != nil {
			return nil, errors.Join(e, err)
		}
		return nil, err
	}
	return r.GetItem(ctx, query.NewMenuItemQuery(item.TenantID).WithID(item.ID).WithPreload())
}

func (r *Repository) UpdateItem(ctx context.Context, dto *dto.UpdateMenuItem) (*entity.Item, error) {
	err := r.dbWithContext(ctx).Transaction(func(tx *gorm.DB) error {
		item, err := r.GetItem(ctx, query.NewMenuItemQuery(dto.TenantID).WithID(dto.ID))
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
		if dto.Category != nil {
			item.Category = *dto.Category
		}
		if dto.FoodType != nil {
			item.FoodType = *dto.FoodType
		}
		if dto.ItemStatus != nil {
			item.ItemStatus = *dto.ItemStatus
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
			if err := r.DeleteFiles(ctx, tx, filesToDelete); err != nil {
				return err
			}

			if err := r.PopulateOwnerID(ctx, tx, item.ID, dto.Files, dto.TenantID); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetItem(ctx, query.NewMenuItemQuery(dto.TenantID).WithID(dto.ID).WithPreload())
}

func (r Repository) GetItem(ctx context.Context, filter *query.MenuItemQuery) (*entity.Item, error) {
	item := &entity.Item{}
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.Item{}).Where("tenant_id = ?", filter.TenantID)
	if filter.ID != nil {
		sess.Where("id = ?", filter.ID)
	}
	if filter.LocationID != nil {
		sess.Where("location_id = ?", filter.LocationID)
	}
	if filter.Orderable != nil {
		sess.Where("item_status = ?", "av")
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

func (r Repository) GetItems(ctx context.Context, filter *query.MenuItemQuery) ([]*entity.Item, error) {
	items := []*entity.Item{}
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.Item{}).Where("tenant_id = ?", filter.TenantID)
	if filter.ID != nil {
		sess.Where("id = ?", filter.ID)
	}
	if filter.LocationID != nil {
		sess.Where("location_id = ?", filter.LocationID)
	}
	if filter.Orderable != nil {
		sess.Where("item_status = ?", "av")
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
