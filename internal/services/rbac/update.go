package rbac

import "grubzo/internal/models/dto"

func (r *RBAC) AddUserRole(payload *dto.AddRole) error {
	return r.repo.CreateRole(
		payload.TenantID,
		payload.Name,
		payload.Permissions,
	)
}

func (r *RBAC) UpdateUserRole(payload *dto.UpdateRoles) error {
	for _, change := range payload.Data {
		if change.Action == 0 {
			if err := r.repo.RemovePermissionsFromRole(payload.TenantID, change.Name, change.Permissions); err != nil {
				return err
			}
		}
		if change.Action == 1 {
			if err := r.repo.AddPermissionsToRole(payload.TenantID, change.Name, change.Permissions); err != nil {
				return err
			}
		}
	}
	return r.Reload()
}
