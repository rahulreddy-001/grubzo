package query

type MenuItemQuery struct {
	TenantID       uint
	ID             *uint
	IDs            []uint
	LocationID     *uint
	Orderable      *bool
	SearchText     *string
	CuisineText    *string
	Limit          *int
	OrderUpdatedAt bool
	Preload        bool
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

func (f *MenuItemQuery) WithIDs(IDs []uint) *MenuItemQuery {
	f.IDs = IDs
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

func (f *MenuItemQuery) WithSearchText(text string) *MenuItemQuery {
	f.SearchText = &text
	return f
}

func (f *MenuItemQuery) WithCuisineText(text string) *MenuItemQuery {
	f.CuisineText = &text
	return f
}

func (f *MenuItemQuery) WithLimit(limit int) *MenuItemQuery {
	f.Limit = &limit
	return f
}

func (f *MenuItemQuery) OrderByUpdatedAtDesc() *MenuItemQuery {
	f.OrderUpdatedAt = true
	return f
}

func (f *MenuItemQuery) WithPreload() *MenuItemQuery {
	f.Preload = true
	return f
}
