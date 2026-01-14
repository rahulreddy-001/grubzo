package query

type OrderQuery struct {
	TenantID   uint
	UserID     *uint
	LocationID *uint
	Status     *string
	PreLoads   bool
}

func NewOrderQuery(tenantID uint) *OrderQuery {
	return &OrderQuery{TenantID: tenantID}
}

func (q *OrderQuery) WithUser(userID uint) *OrderQuery {
	q.UserID = &userID
	return q
}

func (q *OrderQuery) WithLocation(locationID uint) *OrderQuery {
	q.LocationID = &locationID
	return q
}

func (q *OrderQuery) WithStatus(status string) *OrderQuery {
	q.Status = &status
	return q
}

func (q *OrderQuery) WithPreloads() *OrderQuery {
	q.PreLoads = true
	return q
}
