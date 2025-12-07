package query

type TenantUserQuery struct {
	TenantID    uint
	ID          *uint
	Email       *string
	Roles       []string
	LocationID  *uint
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

func (f *TenantUserQuery) WithRole(roles []string) *TenantUserQuery {
	f.Roles = roles
	return f
}

func (f *TenantUserQuery) WithLocationID(id uint) *TenantUserQuery {
	f.LocationID = &id
	return f
}
func (f *TenantUserQuery) WithPreloads() *TenantUserQuery {
	f.WithPreload = true
	return f
}
