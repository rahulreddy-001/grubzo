package rbac

//go:generate go run ../../../cmd/injecttrace -file rabc.go -receiver RBAC -service RBAC
import (
	"context"
	"fmt"
	"grubzo/internal/repository"
	"grubzo/internal/services/rbac/permission"
	"grubzo/internal/services/rbac/role"
	"sync"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type RBAC struct {
	roles      map[uint64]role.Roles
	rolesMutex sync.RWMutex
	repo       *repository.Repository
	db         *gorm.DB
}

func New(repo *repository.Repository) (*RBAC, error) {
	r := &RBAC{
		roles: make(map[uint64]role.Roles),
		repo:  repo,
	}
	if err := r.reload(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to init rbac: %w", err)
	}
	return r, nil
}

func (r *RBAC) IsGranted(tenantID uint64, roleName string, perm permission.Permission) bool {
	r.rolesMutex.RLock()
	defer r.rolesMutex.RUnlock()
	return r.isGranted(tenantID, roleName, perm)
}

func (r *RBAC) isGranted(tenantID uint64, roleName string, perm permission.Permission) bool {
	if roleName == role.Admin {
		return true
	}
	if tenantRoles, ok := r.roles[tenantID]; ok {
		return tenantRoles.HasAndIsGranted(roleName, perm)
	}
	return false
}

func (r *RBAC) IsAllGranted(tenantID uint64, roles []string, perm permission.Permission) bool {
	r.rolesMutex.RLock()
	defer r.rolesMutex.RUnlock()
	for _, roleName := range roles {
		if !r.isGranted(tenantID, roleName, perm) {
			return false
		}
	}
	return true
}

func (r *RBAC) IsAnyGranted(tenantID uint64, roles []string, perm permission.Permission) bool {
	r.rolesMutex.RLock()
	defer r.rolesMutex.RUnlock()
	for _, roleName := range roles {
		if r.isGranted(tenantID, roleName, perm) {
			return true
		}
	}
	return false
}

func (r *RBAC) Reload(ctx context.Context) error {
	ctx, span := otel.Tracer("RBAC").Start(ctx, "RBAC.Reload")
	defer span.End()

	return r.reload(ctx)
}

func (r *RBAC) reload(ctx context.Context) error {
	ctx, span := otel.Tracer("RBAC").Start(ctx, "RBAC.reload")
	defer span.End()

	allRoles, err := r.repo.GetAllUserRoles(ctx, nil)
	if err != nil {
		return err
	}
	tenantRoleMap := make(map[uint64]role.Roles)
	for _, v := range allRoles {
		perms := permission.Permissions{}
		for _, permStr := range v.Permissions {
			perms.Add(permission.Permission(permStr))
		}

		roleObj := role.Role{
			Name:        v.Name,
			Permissions: perms,
		}

		if _, ok := tenantRoleMap[v.TenantID]; !ok {
			tenantRoleMap[v.TenantID] = role.Roles{}
		}
		tenantRoleMap[v.TenantID].Add(roleObj)
	}

	r.rolesMutex.Lock()
	r.roles = tenantRoleMap
	r.rolesMutex.Unlock()
	return nil
}

func (r *RBAC) GetGrantedPermissions(tenantID uint64, roleName string) []permission.Permission {
	if roleName == role.Admin {
		return permission.List
	}

	r.rolesMutex.RLock()
	defer r.rolesMutex.RUnlock()

	if tenantRoles, ok := r.roles[tenantID]; ok {
		if ro, exists := tenantRoles[roleName]; exists {
			return ro.Permissions.Array()
		}
	}
	return nil
}

func (r *RBAC) GetGrantedPermissionsForRoles(tenantID uint64, roleNames []string) []permission.Permission {
	r.rolesMutex.RLock()
	defer r.rolesMutex.RUnlock()
	perms := []permission.Permission{}
	for _, roleName := range roleNames {

		if roleName == role.Admin {
			return permission.List
		}

		if tenantRoles, ok := r.roles[tenantID]; ok {
			if ro, exists := tenantRoles[roleName]; exists {
				perms = append(perms, ro.Permissions.Array()...)
			}
		}
	}
	return perms
}

func (r *RBAC) GetAllPermisssions() []permission.Permission {
	return permission.List
}

func (r *RBAC) GetAllRoles(tenantID uint64) []string {
	r.rolesMutex.RLock()
	defer r.rolesMutex.RUnlock()

	roles := []string{}
	if tenantRoles, ok := r.roles[tenantID]; ok {
		for _, role := range tenantRoles {
			roles = append(roles, role.Name)
		}
	}
	return roles
}

func (r *RBAC) GetAllRolePermissions(ctx context.Context, tenantID uint64) (map[string][]string, error) {
	ctx, span := otel.Tracer("RBAC").Start(ctx, "RBAC.GetAllRolePermissions")
	defer span.End()

	userRoles, err := r.repo.GetAllUserRoles(ctx, &tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles for tenant %d: %w", tenantID, err)
	}
	result := make(map[string][]string)

	for _, ur := range userRoles {
		if len(ur.Permissions) == 0 {
			result[ur.Name] = []string{}
			continue
		}
		perms := make([]string, len(ur.Permissions))
		copy(perms, ur.Permissions)
		result[ur.Name] = perms
	}
	return result, nil
}
