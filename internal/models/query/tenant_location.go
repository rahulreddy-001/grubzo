package query

type TenantLocationQuery struct {
	TenantID          uint
	ID                *uint
	Code              *string
	IsPrimary         *bool
	SearchText        *string
	LocationText      *string
	Limit             *int
	OrderPrimaryFirst bool
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

func (f *TenantLocationQuery) WithSearchText(text string) *TenantLocationQuery {
	f.SearchText = &text
	return f
}

func (f *TenantLocationQuery) WithLocationText(text string) *TenantLocationQuery {
	f.LocationText = &text
	return f
}

func (f *TenantLocationQuery) WithLimit(limit int) *TenantLocationQuery {
	f.Limit = &limit
	return f
}

func (f *TenantLocationQuery) OrderByPrimaryFirst() *TenantLocationQuery {
	f.OrderPrimaryFirst = true
	return f
}
