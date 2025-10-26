package query

type MenuItemQuery struct {
	TenantID   uint
	ID         *uint
	LocationID *uint
	Orderable  *bool
	Preload    bool
}

func NewMenuItemQuery(TenantID uint) *MenuItemQuery {
	return &MenuItemQuery{
		TenantID: TenantID,
	}
}
func (f *MenuItemQuery) WithID(ID uint) *MenuItemQuery {
	f.ID = &ID
	return f
}
func (f *MenuItemQuery) WithLocationID(ID uint) *MenuItemQuery {
	f.LocationID = &ID
	return f
}
func (f *MenuItemQuery) WithOrderable(orderable bool) *MenuItemQuery {
	f.Orderable = &orderable
	return f
}
func (f *MenuItemQuery) WithPreload() *MenuItemQuery {
	f.Preload = true
	return f
}
