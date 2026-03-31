package repository

//go:generate go run ../../cmd/injecttrace -file item.go -receiver Repository -service Repository
import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"slices"
	"strings"
)

type ItemRepository interface {
	CreateItem(ctx context.Context, dto *dto.CreateMenuItem) (*entity.Item, error)
	UpdateItem(ctx context.Context, dto *dto.UpdateMenuItem) (*entity.Item, error)
	GetItem(ctx context.Context, filter *query.MenuItemQuery) (*entity.Item, error)
	GetItems(ctx context.Context, filter *query.MenuItemQuery) ([]*entity.Item, error)
}

func (r Repository) CreateItem(ctx context.Context, dto *dto.CreateMenuItem) (*entity.Item, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.CreateItem")
	defer span.End()

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
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.UpdateItem")
	defer span.End()

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
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.GetItem")
	defer span.End()

	item := &entity.Item{}
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.Item{}).Where("tenant_id = ?", filter.TenantID)
	if filter.ID != nil {
		sess = sess.Where("id = ?", *filter.ID)
	}
	if len(filter.IDs) > 0 {
		sess = sess.Where("id IN ?", filter.IDs)
	}
	if filter.LocationID != nil {
		sess = sess.Where("location_id = ?", *filter.LocationID)
	}
	if filter.Orderable != nil {
		sess = sess.Where("item_status = ?", "av")
	}
	if filter.CuisineText != nil {
		if cuisine := strings.TrimSpace(*filter.CuisineText); cuisine != "" {
			like := "%" + strings.ToLower(cuisine) + "%"
			sess = sess.Where("LOWER(category) LIKE ? OR LOWER(food_type) LIKE ?", like, like)
		}
	}
	if filter.SearchText != nil {
		if text := strings.TrimSpace(*filter.SearchText); text != "" {
			like := "%" + strings.ToLower(text) + "%"
			sess = sess.Where(
				"LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(category) LIKE ?",
				like,
				like,
				like,
			)
		}
	}
	if filter.OrderUpdatedAt {
		sess = sess.Order("updated_at desc")
	}
	if filter.Limit != nil && *filter.Limit > 0 {
		sess = sess.Limit(*filter.Limit)
	}
	if filter.Preload {
		for _, preload := range item.GetPreloads() {
			sess = sess.Preload(preload)
		}
	}
	if err := sess.First(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (r Repository) GetItems(ctx context.Context, filter *query.MenuItemQuery) ([]*entity.Item, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.GetItems")
	defer span.End()

	items := []*entity.Item{}
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.Item{}).Where("tenant_id = ?", filter.TenantID)
	if filter.ID != nil {
		sess = sess.Where("id = ?", *filter.ID)
	}
	if len(filter.IDs) > 0 {
		sess = sess.Where("id IN ?", filter.IDs)
	}
	if filter.LocationID != nil {
		sess = sess.Where("location_id = ?", *filter.LocationID)
	}
	if filter.Orderable != nil {
		sess = sess.Where("item_status = ?", "av")
	}
	if filter.CuisineText != nil {
		if cuisine := strings.TrimSpace(*filter.CuisineText); cuisine != "" {
			like := "%" + strings.ToLower(cuisine) + "%"
			sess = sess.Where("LOWER(category) LIKE ? OR LOWER(food_type) LIKE ?", like, like)
		}
	}
	if filter.SearchText != nil {
		if text := strings.TrimSpace(*filter.SearchText); text != "" {
			like := "%" + strings.ToLower(text) + "%"
			sess = sess.Where(
				"LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(category) LIKE ?",
				like,
				like,
				like,
			)
		}
	}
	if filter.OrderUpdatedAt {
		sess = sess.Order("updated_at desc")
	}
	if filter.Limit != nil && *filter.Limit > 0 {
		sess = sess.Limit(*filter.Limit)
	}
	if filter.Preload {
		for _, preload := range (entity.Item{}).GetPreloads() {
			sess = sess.Preload(preload)
		}
	}
	if err := sess.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
