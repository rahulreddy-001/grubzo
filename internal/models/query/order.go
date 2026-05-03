package query

type OrderQuery struct {
	TenantID     uint
	ID           *uint
	UserID       *uint
	LocationID   *uint
	Status       *string
	Limit        *int
	OrderCreated bool
	PreLoads     bool
}

func NewOrderQuery(tenantID uint) *OrderQuery {
	return &OrderQuery{TenantID: tenantID}
}

func (q *OrderQuery) WithUser(userID uint) *OrderQuery {
	q.UserID = &userID
	return q
}

func (q *OrderQuery) WithID(orderID uint) *OrderQuery {
	q.ID = &orderID
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
