package query

type UserQuery struct {
	TenantID uint64
	ID       *uint64
	Email    *string
	UserId   *uint64
}

func NewUserQuery(tenantID uint64) *UserQuery {
	return &UserQuery{
		TenantID: tenantID,
	}
}

func (f *UserQuery) WithID(id uint64) *UserQuery {
	f.ID = &id
	return f
}

func (f *UserQuery) WithEmail(email string) *UserQuery {
	f.Email = &email
	return f
}

func (f *UserQuery) WithUserId(userId uint64) *UserQuery {
	f.UserId = &userId
	return f
}
