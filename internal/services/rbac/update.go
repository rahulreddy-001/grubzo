package rbac

import (
	"context"
	"grubzo/internal/models/dto"
)

func (r *RBAC) AddUserRole(ctx context.Context, payload *dto.AddRole) error {
	return r.repo.CreateRole(
		ctx,
		payload.TenantID,
		payload.Name,
		payload.Permissions,
	)
}

func (r *RBAC) UpdateUserRole(ctx context.Context, payload *dto.UpdateRoles) error {
	for _, change := range payload.Data {
		if change.Action == 0 {
			if err := r.repo.RemovePermissionsFromRole(ctx, payload.TenantID, change.Name, change.Permissions); err != nil {
				return err
			}
		}
		if change.Action == 1 {
			if err := r.repo.AddPermissionsToRole(ctx, payload.TenantID, change.Name, change.Permissions); err != nil {
				return err
			}
		}
	}
	return r.Reload(ctx)
}
