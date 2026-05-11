package repository

//go:generate go run ../../cmd/injecttrace -file role.go -receiver Repository -service Repository
import (
	"context"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
	"grubzo/internal/models/entity"
	"grubzo/internal/router/ext"
)

type RoleRepository interface {
	GetAllUserRoles(ctx context.Context, tenantID *uint64) ([]entity.UserRole, error)
	CreateRole(ctx context.Context, tenantID uint64, name string, permissions []string) error
	AddPermissionsToRole(ctx context.Context, tenantID uint64, roleName string, permissions []string) error
	RemovePermissionsFromRole(ctx context.Context, tenantID uint64, roleName string, permissions []string) error
}

func (repo *Repository) GetAllUserRoles(ctx context.Context, tenantID *uint64) ([]entity.UserRole, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.GetAllUserRoles")
	defer span.End()

	var userRoles []entity.UserRole
	sess := repo.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.UserRole{})
	if tenantID != nil {
		sess = sess.Where("tenant_id = ?", *tenantID)
	}
	if err := sess.Find(&userRoles).Error; err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (repo *Repository) CreateRole(ctx context.Context, tenantID uint64, name string, permissions []string) error {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.CreateRole")
	defer span.End()

	role := entity.UserRole{
		TenantID:    tenantID,
		Name:        name,
		Permissions: permissions,
	}
	if err := repo.dbWithContext(ctx).Create(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ext.Error(fmt.Sprintf("role %q already exists", name))
		}
		return err
	}
	return nil
}

func (repo *Repository) AddPermissionsToRole(ctx context.Context, tenantID uint64, roleName string, newPerms []string) error {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.AddPermissionsToRole")
	defer span.End()

	var role entity.UserRole
	if err := repo.dbWithContext(ctx).
		Where("tenant_id = ? AND name = ?", tenantID, roleName).
		First(&role).Error; err != nil {
		return err
	}

	existing := make(map[string]struct{})
	for _, p := range role.Permissions {
		existing[p] = struct{}{}
	}
	for _, np := range newPerms {
		if _, exists := existing[np]; !exists {
			role.Permissions = append(role.Permissions, np)
		}
	}

	return repo.dbWithContext(ctx).Model(&role).
		Update("permissions", pq.StringArray(role.Permissions)).Error
}

func (repo *Repository) RemovePermissionsFromRole(ctx context.Context, tenantID uint64, roleName string, permsToRemove []string) error {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.RemovePermissionsFromRole")
	defer span.End()

	var role entity.UserRole
	if err := repo.dbWithContext(ctx).
		Where("tenant_id = ? AND name = ?", tenantID, roleName).
		First(&role).Error; err != nil {
		return err
	}

	toRemove := make(map[string]struct{})
	for _, p := range permsToRemove {
		toRemove[p] = struct{}{}
	}

	newPerms := make([]string, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		if _, remove := toRemove[p]; !remove {
			newPerms = append(newPerms, p)
		}
	}

	return repo.dbWithContext(ctx).Model(&role).
		Update("permissions", pq.StringArray(newPerms)).Error
}

func (repo *Repository) DeleteRole(ctx context.Context, tenantID uint64, roleName string) error {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.DeleteRole")
	defer span.End()

	return repo.dbWithContext(ctx).
		Where("tenant_id = ? AND name = ?", tenantID, roleName).
		Delete(&entity.UserRole{}).
		Error
}
