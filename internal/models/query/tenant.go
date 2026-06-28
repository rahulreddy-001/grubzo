package query

type TenantQuery struct {
	ID        *uint64
	Code      *string
	SubDomain *string
	PreLoads  bool
}

func NewTenantQuery() *TenantQuery {
	return &TenantQuery{
		PreLoads: false,
	}
}

func (q *TenantQuery) WithID(id uint64) *TenantQuery {
	q.ID = &id
	return q
}

func (q *TenantQuery) WithCode(code string) *TenantQuery {
	q.Code = &code
	return q
}

func (q *TenantQuery) WithSubDomain(subDomain string) *TenantQuery {
	q.SubDomain = &subDomain
	return q
}

func (q *TenantQuery) WithPreloads() *TenantQuery {
	q.PreLoads = true
	return q
}
