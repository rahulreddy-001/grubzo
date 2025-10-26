package query

import "grubzo/internal/models/entity"

type TenantUserQuery struct {
	TenantID    uint
	ID          *uint
	Email       *string
	Role        *entity.TenantUserRole
	WithPreload bool
}

func NewTenantUserQuery(tenantID uint) *TenantUserQuery {
	return &TenantUserQuery{
		TenantID: tenantID,
	}
}

func (f *TenantUserQuery) WithID(id uint) *TenantUserQuery {
	f.ID = &id
	return f
}

func (f *TenantUserQuery) WithEmail(email string) *TenantUserQuery {
	f.Email = &email
	return f
}

func (f *TenantUserQuery) WithRole(role entity.TenantUserRole) *TenantUserQuery {
	f.Role = &role
	return f
}
func (f *TenantUserQuery) WithPreloads() *TenantUserQuery {
	f.WithPreload = true
	return f
}
