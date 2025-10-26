package query

type TenantLocationQuery struct {
	TenantID  uint
	ID        *uint
	Code      *string
	IsPrimary *bool
}

func NewTenantLocationQuery(tenantID uint) *TenantLocationQuery {
	return &TenantLocationQuery{
		TenantID: tenantID,
	}
}

func (f *TenantLocationQuery) WithID(id uint) *TenantLocationQuery {
	f.ID = &id
	return f
}

func (f *TenantLocationQuery) WithCode(code string) *TenantLocationQuery {
	f.Code = &code
	return f
}

func (f *TenantLocationQuery) WithPrimary(primary bool) *TenantLocationQuery {
	f.IsPrimary = &primary
	return f
}
