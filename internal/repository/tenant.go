package repository

//go:generate go run ../../cmd/injecttrace -file tenant.go -receiver Repository -service Repository
import (
	"context"
	"errors"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"grubzo/internal/router/ext"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type TenantRepository interface {
	CreateTenant(ctx context.Context, dto *dto.CreateTenant) (*entity.Tenant, error)
	UpdateTenant(ctx context.Context, dto *dto.UpdateTenant) (*entity.Tenant, error)
	GetTenant(ctx context.Context, query *query.TenantQuery) (*entity.Tenant, error)
	GetTenants(ctx context.Context, query *query.TenantQuery) ([]*entity.Tenant, error)
	SaveTenant(ctx context.Context, tenant *entity.Tenant) error
}

func tenantValidator(tenant *entity.Tenant, db *gorm.DB) error {
	// Check for unique name & code
	var existing []entity.Tenant
	model := db.Model(&entity.Tenant{}).Where("code = ? OR name = ?", tenant.Code, tenant.Name)
	if tenant.ID != 0 {
		model = model.Where("id != ?", tenant.ID)
	}
	if err := model.Find(&existing).Error; err != nil {
		return err
	}
	for _, e := range existing {
		if e.Code == tenant.Code {
			return ext.Error("tenant code must be unique")
		}
		if e.Name == tenant.Name {
			return ext.Error("tenant name must be unique")
		}
	}
	return nil
}

func (r *Repository) CreateTenant(ctx context.Context, dto *dto.CreateTenant) (*entity.Tenant, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.CreateTenant")
	defer span.End()

	tenant := &entity.Tenant{
		ID:   *dto.ID,
		Name: dto.Name,
		Code: dto.Code,
	}
	db := r.dbWithContext(ctx)
	if err := tenantValidator(tenant, db); err != nil {
		return nil, err
	}
	if err := db.Create(tenant).Error; err != nil {
		return nil, err
	}
	return tenant, nil
}

func (r *Repository) GetTenant(ctx context.Context, q *query.TenantQuery) (*entity.Tenant, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.GetTenant")
	defer span.End()

	tenant := &entity.Tenant{}
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.Tenant{})

	if q.PreLoads {
		for _, preload := range tenant.GetPreloads() {
			sess = sess.Preload(preload)
		}
	}
	if q.ID != nil {
		sess = sess.Where("id = ?", *q.ID)
	}
	if q.Code != nil {
		sess = sess.Where("code = ?", *q.Code)
	}

	if err := sess.First(tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ext.Error("Tenant not found")
		}
		return nil, err
	}
	return tenant, nil
}

func (r *Repository) GetTenants(ctx context.Context, q *query.TenantQuery) ([]*entity.Tenant, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.GetTenants")
	defer span.End()

	var tenants []*entity.Tenant
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.Tenant{})

	if q.PreLoads {
		for _, preload := range (entity.Tenant{}).GetPreloads() {
			sess = sess.Preload(preload)
		}
	}
	if q.ID != nil {
		sess = sess.Where("id = ?", *q.ID)
	}
	if q.Code != nil {
		sess = sess.Where("code = ?", *q.Code)
	}

	if err := sess.Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

func (r *Repository) UpdateTenant(ctx context.Context, dto *dto.UpdateTenant) (*entity.Tenant, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.UpdateTenant")
	defer span.End()

	updates := map[string]any{}
	if dto.Name != nil {
		updates["name"] = *dto.Name
	}
	if dto.Code != nil {
		updates["code"] = *dto.Code
	}

	if len(updates) == 0 {
		return r.GetTenant(ctx, query.NewTenantQuery().WithID(dto.ID).WithPreloads())
	}

	if err := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.Tenant{}).
		Where("id = ?", dto.ID).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	return r.GetTenant(ctx, query.NewTenantQuery().WithID(dto.ID).WithPreloads())
}

func (r *Repository) SaveTenant(ctx context.Context, tenant *entity.Tenant) error {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.SaveTenant")
	defer span.End()

	db := r.dbWithContext(ctx)
	if err := tenantValidator(tenant, db); err != nil {
		return err
	}
	return db.Save(tenant).Error
}
