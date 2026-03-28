package rbac

//go:generate go run ../../../cmd/injecttrace -file update.go -receiver RBAC -service RBAC
import (
	"context"
	"go.opentelemetry.io/otel"
	"grubzo/internal/models/dto"
)

func (r *RBAC) AddUserRole(ctx context.Context, payload *dto.AddRole) error {
	ctx, span := otel.Tracer("RBAC").Start(ctx, "RBAC.AddUserRole")
	defer span.End()

	return r.repo.CreateRole(
		ctx,
		payload.TenantID,
		payload.Name,
		payload.Permissions,
	)
}

func (r *RBAC) UpdateUserRole(ctx context.Context, payload *dto.UpdateRoles) error {
	ctx, span := otel.Tracer("RBAC").Start(ctx, "RBAC.UpdateUserRole")
	defer span.End()

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
