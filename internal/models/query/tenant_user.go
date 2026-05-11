package query

type TenantUserQuery struct {
	TenantID    uint64
	ID          *uint64
	Email       *string
	Roles       []string
	LocationID  *uint64
	WithPreload bool
}

func NewTenantUserQuery(tenantID uint64) *TenantUserQuery {
	return &TenantUserQuery{
		TenantID: tenantID,
	}
}

func (f *TenantUserQuery) WithID(id uint64) *TenantUserQuery {
	f.ID = &id
	return f
}

func (f *TenantUserQuery) WithEmail(email string) *TenantUserQuery {
	f.Email = &email
	return f
}

func (f *TenantUserQuery) WithRole(roles []string) *TenantUserQuery {
	f.Roles = roles
	return f
}

func (f *TenantUserQuery) WithLocationID(id uint64) *TenantUserQuery {
	f.LocationID = &id
	return f
}
func (f *TenantUserQuery) WithPreloads() *TenantUserQuery {
	f.WithPreload = true
	return f
}
