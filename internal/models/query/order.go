package query

type OrderQuery struct {
	TenantID     uint64
	ID           *uint64
	UserID       *uint64
	LocationID   *uint64
	Status       *string
	Limit        *int
	OrderCreated bool
	PreLoads     bool
}

func NewOrderQuery(tenantID uint64) *OrderQuery {
	return &OrderQuery{TenantID: tenantID}
}

func (q *OrderQuery) WithUser(userID uint64) *OrderQuery {
	q.UserID = &userID
	return q
}

func (q *OrderQuery) WithID(orderID uint64) *OrderQuery {
	q.ID = &orderID
	return q
}

func (q *OrderQuery) WithLocation(locationID uint64) *OrderQuery {
	q.LocationID = &locationID
	return q
}

func (q *OrderQuery) WithStatus(status string) *OrderQuery {
	q.Status = &status
	return q
}

func (q *OrderQuery) WithLimit(limit int) *OrderQuery {
	q.Limit = &limit
	return q
}

func (q *OrderQuery) OrderByCreatedAtDesc() *OrderQuery {
	q.OrderCreated = true
	return q
}

func (q *OrderQuery) WithPreloads() *OrderQuery {
	q.PreLoads = true
	return q
}
