package query

type UserQuery struct {
	TenantID uint
	ID       *uint
	Email    *string
	UserId   *uint
}

func NewUserQuery(tenantID uint) *UserQuery {
	return &UserQuery{
		TenantID: tenantID,
	}
}

func (f *UserQuery) WithID(id uint) *UserQuery {
	f.ID = &id
	return f
}

func (f *UserQuery) WithEmail(email string) *UserQuery {
	f.Email = &email
	return f
}

func (f *UserQuery) WithUserId(userId uint) *UserQuery {
	f.UserId = &userId
	return f
}
