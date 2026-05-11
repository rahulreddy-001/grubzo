package actor

type Actor struct {
	UserID     uint64   `json:"user_id"`
	TenantID   uint64   `json:"tenant_id"`
	LocationID uint64   `json:"location_id"`
	Email      string   `json:"email"`
	Type       string   `json:"type"`
	Roles      []string `json:"roles"`
}
