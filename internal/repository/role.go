package repository

import (
	"errors"
	"fmt"
	"grubzo/internal/models/entity"
	"grubzo/internal/router/ext"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type RoleRepository interface {
	GetAllUserRoles(tenantID *uint) ([]entity.UserRole, error)
	CreateRole(tenantID uint, name string, permissions []string) error
	AddPermissionsToRole(tenantID uint, roleName string, permissions []string) error
	RemovePermissionsFromRole(tenantID uint, roleName string, permissions []string) error
}

func (repo *Repository) GetAllUserRoles(tenantID *uint) ([]entity.UserRole, error) {
	var userRoles []entity.UserRole
	sess := repo.db.Session(&gorm.Session{}).Model(&entity.UserRole{})
	if tenantID != nil {
		sess = sess.Where("tenant_id = ?", *tenantID)
	}
	if err := sess.Find(&userRoles).Error; err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (repo *Repository) CreateRole(tenantID uint, name string, permissions []string) error {
	role := entity.UserRole{
		TenantID:    tenantID,
		Name:        name,
		Permissions: permissions,
	}
	if err := repo.db.Create(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ext.Error(fmt.Sprintf("role %q already exists", name))
		}
		return err
	}
	return nil
}

func (repo *Repository) AddPermissionsToRole(tenantID uint, roleName string, newPerms []string) error {
	var role entity.UserRole
	if err := repo.db.
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

	return repo.db.Model(&role).
		Update("permissions", pq.StringArray(role.Permissions)).Error
}

func (repo *Repository) RemovePermissionsFromRole(tenantID uint, roleName string, permsToRemove []string) error {
	var role entity.UserRole
	if err := repo.db.
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

	return repo.db.Model(&role).
		Update("permissions", pq.StringArray(newPerms)).Error
}

func (repo *Repository) DeleteRole(tenantID uint, roleName string) error {
	return repo.db.
		Where("tenant_id = ? AND name = ?", tenantID, roleName).
		Delete(&entity.UserRole{}).
		Error
}
